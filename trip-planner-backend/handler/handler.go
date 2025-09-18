package handler

import (
	findDestinationController "trip-planner/modules/tripPlanner/controller/findDestination"
	generateItineraryController "trip-planner/modules/tripPlanner/controller/generateItinerary"
	processItineraryController "trip-planner/modules/tripPlanner/controller/processItinerary"

	"github.com/gin-gonic/gin"
)

func RouteRequests(ginServer *gin.Engine) {
	apiGroup := ginServer.Group("tripPlanner")
	apiGroup.POST("/findDestination", findDestinationController.FindDestination)
	apiGroup.POST("/findDestination/", findDestinationController.FindDestination)
	apiGroup.POST("/generateItinerary", generateItineraryController.GenerateItinerary)
	apiGroup.POST("/generateItinerary/", generateItineraryController.GenerateItinerary)
	apiGroup.POST("/processItinerary", processItineraryController.ProcessItinerary)
	apiGroup.POST("/processItinerary/", processItineraryController.ProcessItinerary)
}
