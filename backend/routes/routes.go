package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shjung-dev/1to1Chat/backend/controllers"
	"github.com/shjung-dev/1to1Chat/backend/middleware"
)

func SetUpRoutes(r *gin.Engine) {
	r.POST("/signup", controllers.Signup())
	r.POST("/login", controllers.Login())
	r.POST("/refresh", controllers.RefreshTokenHandler())

	protected := r.Group("/")

	protected.Use(middleware.Authenticate())
	{
		protected.GET("/users", controllers.GetUsers())
		protected.GET("/user/:username", controllers.SearchUser())
	}

}
