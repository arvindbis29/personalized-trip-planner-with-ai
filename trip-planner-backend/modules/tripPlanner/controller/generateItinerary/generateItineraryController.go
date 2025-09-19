package generateItineraryController

import (
	generateItineraryModel "trip-planner-backend/modules/tripPlanner/model/generateItinerary"

	"github.com/gin-gonic/gin"
)

func GenerateItinerary(ginCtx *gin.Context) {
	ginCtx.Status(200)
}

func BindInputParams(ginCtx *gin.Context) (InputParams generateItineraryModel.ApiInputParams, err error) {
	bindErr := ginCtx.ShouldBindBodyWithJSON(&InputParams)
	return InputParams, bindErr
}

func ReturnApiResponse(ginCtx *gin.Context, apiCode int, apiResponse generateItineraryModel.ApiResponse) {
	ginCtx.JSON(apiCode, apiResponse)
}
