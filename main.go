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

	user.MigrateUser(config.DB)
	activity.MigrateActivity(config.DB)

	userRepository := user.NewRepository(config.DB)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	authService := auth.NewService(userRepository)
	authHandler := auth.NewHandler(authService)

	r := gin.Default()
	routes.AuthRoutes(r, authHandler)
	routes.UserRoutes(r, userHandler)

	r.Run(":8080")
}