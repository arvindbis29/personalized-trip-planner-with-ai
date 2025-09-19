package findDestinationModel

import (
	"trip-planner-backend/utilities/globalFunctions"

	"github.com/gin-gonic/gin"
)

type BindingInputParams struct {
	UserId                int    `json:"user_id" binding:"required"`
	IsInternationalTravel bool   `json:"is_international_travel"`
	TravelDays            int    `json:"travel_days"`
	TravelDateTime        string `json:"travel_date_time"`
	TravelingMethod       string `json:"traveling_method"`
	TripNature            string `json:"trip_nature"`
	PersonCount           int    `json:"person_count"`
	GroupDemographic      string `json:"group_demographic"`
	Budget                string `json:"budget"`
	CustomRequirement     string `json:"custom_requirement"`
	PreferredLocation     string `json:"preferred_location"`
}

type ApiResponse struct {
	Code     int            `json:"code"`
	Status   string         `json:"status"`
	Response map[string]any `json:"response"`
}

func CreateApplicationLogs(ginCtx *gin.Context, apiInputParams BindingInputParams, apiResponse ApiResponse) {

	fileName := "find_destination"

	logData := map[string]any{}
	logData["user_id"] = globalFunctions.ConvertValueToString(apiInputParams.UserId)
	logData["is_international_travel"] = globalFunctions.ConvertValueToString(apiInputParams.IsInternationalTravel)
	logData["travel_days"] = globalFunctions.ConvertValueToString(apiInputParams.TravelDays)
	logData["travel_date_time"] = globalFunctions.ConvertValueToString(apiInputParams.TravelDateTime)
	logData["traveling_method"] = globalFunctions.ConvertValueToString(apiInputParams.TravelingMethod)
	logData["trip_nature"] = globalFunctions.ConvertValueToString(apiInputParams.TripNature)
	logData["person_count"] = globalFunctions.ConvertValueToString(apiInputParams.PersonCount)
	logData["group_demographic"] = globalFunctions.ConvertValueToString(apiInputParams.GroupDemographic)
	logData["budget"] = globalFunctions.ConvertValueToString(apiInputParams.Budget)
	logData["custom_requirement"] = globalFunctions.ConvertValueToString(apiInputParams.CustomRequirement)
	logData["preferred_location"] = globalFunctions.ConvertValueToString(apiInputParams.PreferredLocation)

	logData["code"] = globalFunctions.ConvertValueToString(apiResponse.Code)
	logData["status"] = globalFunctions.ConvertValueToString(apiResponse.Status)
	logData["response"] = globalFunctions.ConvertJsonValToString(apiResponse.Response)

	globalFunctions.WriteJsonLogs(ginCtx, fileName, logData)
}
