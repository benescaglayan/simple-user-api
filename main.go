package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"user-service/controller"
	"user-service/repository"
	"user-service/service"
)

func main() {
	router := gin.Default()

	validator := validator.New()
	database := repository.InitDatabase(os.Getenv("MONGO_URI"))
	userRepository := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService, validator)

	router.GET("/users", userController.GetAll)
	router.GET("/users/:id", userController.GetById)
	router.POST("/users", userController.Create)
	router.PATCH("/users/:id", userController.UpdateById)
	router.DELETE("/users/:id", userController.DeleteById)

	err := router.Run()
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
