package findDestinationController

import (
	"net/http"
	findDestinationModel "trip-planner-backend/modules/tripPlanner/model/findDestination"
	"trip-planner-backend/utilities/globalFunctions"

	"github.com/gin-gonic/gin"
)

func FindDestination(ginCtx *gin.Context) {
	apiInputParam, bindErr := BindInputParams(ginCtx)
	apiResponse := findDestinationModel.ApiResponse{}
	if bindErr != nil {
		apiResponse.Code= http.StatusBadRequest
		apiResponse.Status = "Failure"
		apiResponse.Response = map[string]any{
			"reponse" : nil,
			"error" : bindErr.Error(),
		}
		ReturnApiResponse(ginCtx, http.StatusBadRequest, apiResponse)
	}
	globalFunctions.WriteJsonLogs(ginCtx, "arvind_testing", map[string]any{"arvind": "k"})
	ginCtx.Status(200)
}

func BindInputParams(ginCtx *gin.Context) (InputParams findDestinationModel.BindingInputParams, err error) {
	bindErr := ginCtx.ShouldBindBodyWithJSON(&InputParams)
	return InputParams, bindErr
}

func ReturnApiResponse(ginCtx *gin.Context, apiCode int, apiResponse findDestinationModel.ApiResponse) {
	ginCtx.JSON(apiCode, apiResponse)
}