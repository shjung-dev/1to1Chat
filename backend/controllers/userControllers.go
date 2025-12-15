package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shjung-dev/1to1Chat/backend/config"
	"github.com/shjung-dev/1to1Chat/backend/helpers"
	"github.com/shjung-dev/1to1Chat/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

var userCollection = config.OpenCollection("users")

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		defer cancel()

		var user models.User

		//Get user input
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Validate user input
		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		//ensure every username is unique
		count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}

		user.Password = helpers.HashPassword(user.Password)
		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		accessToken, refreshToken := helpers.GenerateToken(user.User_id, *user.Username)
		user.Token = &accessToken
		user.Refresh_token = &refreshToken

		_, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username. User is not found"})
		}

		passwordIsValid, msg := helpers.VerifyPassword(*foundUser.Password, *user.Password)

		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		}

		token, refreshToken := helpers.GenerateToken(foundUser.User_id, *foundUser.Username)
		helpers.UpdateAllToken(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, gin.H{
			"user":          foundUser,
			"token":         token,
			"refresh_token": refreshToken,
		})
	}
}

func RefreshTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//Browser sends one field only -> refresh_token from local storage
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
			return
		}

		//Validate Refresh Token
		claims, err := helpers.ValidateToken(body.RefreshToken)
		if err != nil {
			//Refresh token is invalid or expired -> force user to login back
			c.JSON(http.StatusUnauthorized, gin.H{"error": "relogin"})
			return
		}

		userID := claims.UserID

		var user models.User

		err = userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)

		if err != nil {
			//User not found
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not found"})
		}

		//Generate new tokens
		newAccessToken, newRefreshToken := helpers.GenerateToken(user.User_id, *user.Username)
		helpers.UpdateAllToken(newAccessToken, newRefreshToken, user.User_id)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		})
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		

		/*
		-> cursor is of type *mongo.Cursor.
		-> It represents a pointer to the result set of your query.
		-> Unlike FindOne, which returns a single document, Find can return multiple documents. 
		MongoDB doesn’t load all documents at once; it gives you a cursor to iterate over the results efficiently.
		*/
		cursor , err := userCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError , gin.H{"error":err.Error()})
			return
		}

		defer cursor.Close(ctx)

		var users []models.User

		if err := cursor.All(ctx , &users); err != nil {
			c.JSON(http.StatusInternalServerError , gin.H{"error":err.Error()})
			return
		}

		c.JSON(http.StatusOK , users)
	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		requestUserID := c.Param("id")
		

		var user models.User

		err := userCollection.FindOne(ctx , bson.M{"user_id":requestUserID}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}


		c.JSON(http.StatusOK, user)
	}
}
