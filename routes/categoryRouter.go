package routes

import (
	controller "Golang-Management-Platform/controllers"

	"github.com/gin-gonic/gin"
)

// used to be menu

func CategoryRoutes(incomingRoutes *gin.Engine) {
	// Define routes to get, create, and update categories
	incomingRoutes.GET("/categories", controller.GetCategories())
	incomingRoutes.GET("/categories/:category_id", controller.GetCategory())
	incomingRoutes.POST("/categories", controller.CreateCategory())
	incomingRoutes.PATCH("/categories/:category_id", controller.UpdateCategory())
}
