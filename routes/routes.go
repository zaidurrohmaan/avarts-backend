package routes

import (
	"avarts/controllers"
	"avarts/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authCtrl *controllers.AuthController) {
	r.POST("auth/google-login", authCtrl.GoogleLogin)

	protected := r.Group("/").Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", authCtrl.Profile)
	}
}