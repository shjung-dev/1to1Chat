package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shjung-dev/1to1Chat/backend/config"
	"github.com/shjung-dev/1to1Chat/backend/helpers"
	"github.com/shjung-dev/1to1Chat/backend/routes"
)

func main(){
	err := godotenv.Load()

	if err != nil{
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	
	//Connect to the database
	config.ConnectDatabase()
	
	
	Key := config.GenerateRandomKey()
	helpers.SetJWTKey(Key)
	
	//Create the router using gin
	r := gin.Default()
	routes.SetUpRoutes(r)

	//Start the server using gin
	r.Run(":" + port)
	log.Println("Server is running on localhost:8080")


}
