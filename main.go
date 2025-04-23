package main

import (
	"avarts/activity"
	"avarts/auth"
	"avarts/config"
	"avarts/routes"
	"avarts/user"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitDB()
	config.InitAWS()

	user.MigrateUser(config.DB)
	activity.MigrateActivity(config.DB)

	userRepository := user.NewRepository(config.DB)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	authService := auth.NewService(userRepository)
	authHandler := auth.NewHandler(authService)

	activityRepository := activity.NewRepository(config.DB)
	activityService := activity.NewService(activityRepository)
	activityHandler := activity.NewHandler(activityService)

	r := gin.Default()
	
	api := r.Group("/api")
	routes.AuthRoutes(api, authHandler)
	routes.UserRoutes(api, userHandler)
	routes.ActivityRoutes(api, activityHandler)
	routes.UploadRoutes(api, userHandler, activityHandler)

	r.Run(":8080")
}
