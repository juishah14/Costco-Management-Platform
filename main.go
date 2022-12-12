package main

// get something from routes, goes to your controllers, goes to ur models and from the models u have a database connection

import (
	"os"

	middleware "Golang-Management-Platform/middleware"
	routes "Golang-Management-Platform/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// May need to load .env file here

	// Set port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Using Gin to create a router for us
	router := gin.New()
	router.Use(gin.Logger())

	// Declare routes
	routes.AuthRoutes(router)

	// Using middleware to ensure that all following routes are protected, meaning that a user
	// cannot access these routes without first logging in/having a token
	router.Use(middleware.Authentication())

	routes.UserRoutes(router)
	routes.ProductRoutes(router)
	routes.CategoryRoutes(router)
	routes.MembershipRoutes(router)
	routes.AccountRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port)
}
