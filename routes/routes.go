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
	users := r.Group("/users", middlewares.AuthMiddleware())
	{
		users.GET("/:username", userHandler.GetUser)
		users.GET("/me", userHandler.GetMyProfile)
		users.PATCH("/me", userHandler.UpdateUser)
		users.DELETE("/me", userHandler.DeleteUser)
	}
}

func ActivityRoutes(r *gin.RouterGroup, activityHandler *activity.Handler) {
	activities := r.Group("/activities", middlewares.AuthMiddleware())
	{
		// CRUD Activity
		activities.POST("", activityHandler.CreateActivity)
		activities.GET("", activityHandler.GetAllActivities)
		activities.GET("/:id", activityHandler.GetActivity)
		activities.DELETE("/:id", activityHandler.DeleteActivity)

		// Likes
		activities.POST("/:id/like", activityHandler.CreateLike)
		activities.DELETE("/:id/like", activityHandler.DeleteLike)

		// Comments
		activities.POST("/:id/comments", activityHandler.CreateComment)
		activities.DELETE("/comments/:id", activityHandler.DeleteComment)
	}
}

func UploadRoutes(r *gin.RouterGroup, userHandler *user.Handler, activityHandler *activity.Handler) {
	protected := r.Group("/uploads", middlewares.AuthMiddleware())
	{
		protected.POST("/avatar", userHandler.UploadAvatar)
		protected.POST("/activity-photo", activityHandler.UploadActivityPhoto)
	}
}
