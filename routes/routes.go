package routes

import (
	"avarts/activity"
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

func ActivityRoutes(r *gin.Engine, activityHandler *activity.Handler) {
	protected := r.Group("/").Use(middlewares.AuthMiddleware())
	{
		protected.POST("/photos", activityHandler.UploadActivityPhoto)
		protected.POST("/activities", activityHandler.PostActivity)
		protected.GET("/activities/:id", activityHandler.GetActivityByID)
		protected.GET("/activities", activityHandler.GetAllActivities)
	}
}