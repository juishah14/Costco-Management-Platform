package routes

import (
	controller "Golang-Management-Platform/controllers"

	"github.com/gin-gonic/gin"
)

// used to be table

func MembershipRoutes(incomingRoutes *gin.Engine) {
	// Define routes to get, create, and update memberships
	incomingRoutes.GET("/memberships", controller.GetMemberships())
	incomingRoutes.GET("/memberships/:membership_id", controller.GetMembership())
	incomingRoutes.POST("/memberships", controller.CreateMembership())
	incomingRoutes.PATCH("/memberships/:membership_id", controller.UpdateMembership())
}
