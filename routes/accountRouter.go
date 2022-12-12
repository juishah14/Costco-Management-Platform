package routes

import (
	controller "Golang-Management-Platform/controllers"

	"github.com/gin-gonic/gin"
)

// Need to change to accounts

func AccountRoutes(incomingRoutes *gin.Engine) {
	// Define routes to get, create, and update accounts
	// Note that accounts belong to a specific membership
	incomingRoutes.GET("/accounts", controller.GetAccounts())
	incomingRoutes.GET("/accounts/:account_id", controller.GetAccount())
	incomingRoutes.GET("/accounts-membership/:membership_id", controller.GetAccountsByMembership())
	incomingRoutes.POST("/accounts", controller.CreateAccount())
	incomingRoutes.PATCH("/accounts/:account_id", controller.UpdateAccount())
}
