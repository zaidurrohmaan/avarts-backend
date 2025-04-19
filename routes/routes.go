package routes

import (
	"avarts/auth"
	"avarts/middlewares"
	"avarts/user"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, authHandler *auth.Handler) {
	r.POST("auth/google-login", authHandler.GoogleLogin)
}

func UserRoutes(r *gin.Engine, userHandler *user.Handler) {
	protected := r.Group("/").Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile/:username", userHandler.Profile)
		protected.GET("/profile/me", userHandler.MyProfile)
		protected.PATCH("/profile/update/", userHandler.UpdateProfile)
	}
}
