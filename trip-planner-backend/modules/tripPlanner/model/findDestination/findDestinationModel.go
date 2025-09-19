package findDestinationModel

import "github.com/gin-gonic/gin"

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

func CreateApplicationLogs(ginCtx *gin.Context, apiInputParams BindingInputParams) {
}