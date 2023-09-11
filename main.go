package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jenniekibiri/jwt-go/controllers"
	"github.com/jenniekibiri/jwt-go/initializers"
	"github.com/jenniekibiri/jwt-go/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDb()
}

func main() {

	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run()

}
