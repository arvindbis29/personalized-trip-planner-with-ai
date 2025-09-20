package generateItineraryController

import (
	"net/http"
	generateItineraryModel "trip-planner-backend/modules/tripPlanner/model/generateItinerary"

	"github.com/gin-gonic/gin"
)

func GenerateItinerary(ginCtx *gin.Context) {
	apiInputParam, bindErr := BindInputParams(ginCtx)
	apiResponse := generateItineraryModel.ApiResponse{}
	defer func() {
		generateItineraryModel.CreateApplicationLogs(ginCtx, apiInputParam, apiResponse)
	}()

	if bindErr != nil {
		apiResponse.Code = http.StatusBadRequest
		apiResponse.Status = "Failure"
		apiResponse.Error = bindErr.Error()
		ReturnApiResponse(ginCtx, http.StatusBadRequest, apiResponse)
	}

	userQuery := generateItineraryModel.GenerateUserQuery(apiInputParam)
	resp, respErr := generateItineraryModel.GenerateItinerarySuggestion(ginCtx, userQuery)
	if respErr != nil {
		apiResponse.Code = http.StatusInternalServerError
		apiResponse.Status = "Success"
		apiResponse.Error = respErr.Error()
	}
	apiResponse.Code = http.StatusOK
	apiResponse.Status = "Success"
	apiResponse.Response = resp
	ReturnApiResponse(ginCtx, http.StatusOK, apiResponse)
}

func BindInputParams(ginCtx *gin.Context) (InputParams generateItineraryModel.ApiInputParams, err error) {
	bindErr := ginCtx.ShouldBindBodyWithJSON(&InputParams)
	return InputParams, bindErr
}

func ReturnApiResponse(ginCtx *gin.Context, apiCode int, apiResponse generateItineraryModel.ApiResponse) {
	ginCtx.JSON(apiCode, apiResponse)
}
