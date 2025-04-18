package main

import (
	"avarts/config"
	"avarts/controllers"
	"avarts/models"
	"avarts/repository"
	"avarts/routes"
	"avarts/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	models.MigrateUser(config.DB)

	repo := repository.NewUserRepository(config.DB)
	service := services.NewAuthService(repo)
	ctrl := controllers.NewAuthController(service)

	r := gin.Default()
	routes.RegisterRoutes(r, ctrl)

	r.Run(":8080")
}