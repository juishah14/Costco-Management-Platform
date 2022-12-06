package routes

import (
	controller "Golang-Management-Platform/controllers"

	"github.com/gin-gonic/gin"
)

// used to be food

func ProductRoutes(incomingRoutes *gin.Engine) {
	// Define routes to get, create, and update food items
	incomingRoutes.GET("/products", controller.GetProducts())
	incomingRoutes.GET("/products/:product_id", controller.GetProduct())
	incomingRoutes.POST("/products", controller.CreateProduct())
	incomingRoutes.PATCH("/products/:product_id", controller.UpdateProduct())
}
