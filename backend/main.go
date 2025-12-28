package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shjung-dev/1to1Chat/backend/config"
	"github.com/shjung-dev/1to1Chat/backend/helpers"
	"github.com/shjung-dev/1to1Chat/backend/network"
	"github.com/shjung-dev/1to1Chat/backend/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	jwtKey := os.Getenv("JWT_SECRET")

	// Connect to the database
	config.ConnectDatabase()

	// Set JWT key for helper functions
	helpers.SetJWTKey(jwtKey)

	// Create Gin router
	r := gin.Default()

	// Enable CORS for frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ---------------- WebSocket ----------------
	r.GET("/ws", func(c *gin.Context) {   
		token := c.Query("token")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token or recipient"})
			return
		}
		
		// Verify JWT token
		claims, err := helpers.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		
		username := claims.Username

		//Create one room per user
		roomID := network.UserRoom(username)
		chatRoom := network.GetRoom(roomID)
		   
		req := c.Request //c.Request is the underlying *http.Request object that Gin wraps.
		/*
			&http.Request{
			    Method: "GET",
			    URL: &url.URL{Scheme:"http", Host:"localhost:8080", Path:"/ws", RawQuery:"token=abc123&"},
			    Header: http.Header{"Origin":[]string{"http://localhost:3000"}},
			    Body: nil,
			    ...
			}
		*/
		q := req.URL.Query() //req.URL.Query() returns a url.Values type
		q.Set("username", username)
		/* Example (Add key value pair of username with sender -> "alice" as an example):
		q is still a url.Values map:
			map[string][]string{
				"token": {"abc123"},
				"username": {"alice"},
			}
		*/
		req.URL.RawQuery = q.Encode() //q.Encode() converts the url.Values map into a URL-encoded query string
		//map into -> "token=abc123&username=alice"

		chatRoom.ServeHTTP(c.Writer, req)
	})

	// ---------------- Other API routes ----------------
	routes.SetUpRoutes(r)

	// Start the server
	log.Println("Server is running on localhost:" + port)
	r.Run(":" + port)
}
