package findDestinationController

import (
	"net/http"
	findDestinationModel "trip-planner-backend/modules/tripPlanner/model/findDestination"

	"github.com/gin-gonic/gin"
)

func FindDestination(ginCtx *gin.Context) {
	apiInputParam, bindErr := BindInputParams(ginCtx)
	apiResponse := findDestinationModel.ApiResponse{}
	defer func() {
		findDestinationModel.CreateApplicationLogs(ginCtx, apiInputParam, apiResponse)
	}()

	if bindErr != nil {
		apiResponse.Code = http.StatusBadRequest
		apiResponse.Status = "Failure"
		apiResponse.Response = map[string]any{
			"reponse": nil,
			"error":   bindErr.Error(),
		}
		ReturnApiResponse(ginCtx, http.StatusBadRequest, apiResponse)
	}

	ReturnApiResponse(ginCtx, http.StatusOK, apiResponse)
}

func BindInputParams(ginCtx *gin.Context) (InputParams findDestinationModel.BindingInputParams, err error) {
	bindErr := ginCtx.ShouldBindBodyWithJSON(&InputParams)
	return InputParams, bindErr
}

func ReturnApiResponse(ginCtx *gin.Context, apiCode int, apiResponse findDestinationModel.ApiResponse) {
	ginCtx.JSON(apiCode, apiResponse)
}
