package routes

import (
	"avarts/activity"
	"avarts/auth"
	"avarts/middlewares"
	"avarts/user"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, authHandler *auth.Handler) {
	r.POST("auth/google-login", authHandler.GoogleLogin)
}

func UserRoutes(r *gin.RouterGroup, userHandler *user.Handler) {
	protected := r.Group("/profile", middlewares.AuthMiddleware())
	{
		protected.GET("/:username", userHandler.Profile)
		protected.GET("/me", userHandler.MyProfile)
		protected.PATCH("/update", userHandler.UpdateProfile)
		protected.DELETE("", userHandler.DeleteActivity)
	}
}

func ActivityRoutes(r *gin.RouterGroup, activityHandler *activity.Handler) {
	protected := r.Group("/", middlewares.AuthMiddleware())
	{
		activities := protected.Group("/activities")
		activities.POST("", activityHandler.PostActivity)
		activities.GET("/:id", activityHandler.GetActivityByID)
		activities.GET("", activityHandler.GetAllActivities)
		activities.DELETE("", activityHandler.DeleteActivity)

		likes := protected.Group("/like")
		likes.POST("", activityHandler.CreateLike)
		likes.DELETE("", activityHandler.DeleteLike)

		comments := protected.Group("/comment")
		comments.POST("", activityHandler.CreateComment)
		comments.DELETE("", activityHandler.DeleteComment)

		protected.POST("/photos", activityHandler.UploadActivityPhoto)
	}
}

func UploadRoutes(r *gin.RouterGroup, userHandler *user.Handler, activityHandler *activity.Handler) {
	protected := r.Group("/upload", middlewares.AuthMiddleware())
	{
		protected.POST("/avatar", userHandler.UploadAvatar)
		protected.POST("/activity", activityHandler.UploadActivityPhoto)
	}
}
