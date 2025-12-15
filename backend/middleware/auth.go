package middleware

import (
	"errors"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/shjung-dev/1to1Chat/backend/helpers"
)


func Authenticate() gin.HandlerFunc{

	return func(c *gin.Context){
		//Header usually in the form of -> Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == ""{
			c.JSON(http.StatusUnauthorized , gin.H{"error" : "authorization header is required"})
			c.Abort()
			return
		}

		//Remove Bearer because we are only interested in the token
		authHeader = strings.TrimPrefix(authHeader , "Bearer")
		
		claims , err := helpers.ValidateToken(authHeader)

		if err != nil {
			if errors.Is(err , helpers.ErrTokenExpired){
				//Access token expired - client should call /refresh
				c.JSON(401 , gin.H{"error": "access token expired"})
				c.Abort()
				return
			}
			c.JSON(401 , gin.H{"error":"invalid token"})
			c.Abort()
			return
		}

		//Token valid
		c.Set("claims" , claims)
		c.Next() 

	}
}