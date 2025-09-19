package findDestinationModel

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type BindingInputParams struct {
	UserId                int    `json:"user_id" binding:"required"`
	UserLocation          string `json:"user_location"`
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

func GenerateUserQuery(apiInputParams BindingInputParams) (string) {
	var b strings.Builder

	b.WriteString("Given the user inputs below, suggest suitable trip locations tailored to their preferences and constraints.\n\n")
	b.WriteString("User inputs:\n")

	if apiInputParams.UserLocation == "" {
		apiInputParams.UserLocation = "India"
	}
	b.WriteString(fmt.Sprintf("- User is located in %s\n", apiInputParams.UserLocation))

	// Travel type (include explicit yes/no for clarity)
	if apiInputParams.IsInternationalTravel {
		b.WriteString("- International travel: yes\n")
	} else {
		b.WriteString("- International travel: no (domestic)\n")
	}

	// Add other fields only when they are non-zero / non-empty
	if apiInputParams.TravelDays > 0 {
		b.WriteString(fmt.Sprintf("- Travel duration (days): %d\n", apiInputParams.TravelDays))
	}
	if strings.TrimSpace(apiInputParams.TravelDateTime) != "" {
		b.WriteString(fmt.Sprintf("- Travel dates / time: %s\n", strings.TrimSpace(apiInputParams.TravelDateTime)))
	}
	if strings.TrimSpace(apiInputParams.TravelingMethod) != "" {
		b.WriteString(fmt.Sprintf("- Preferred traveling method: %s\n", strings.TrimSpace(apiInputParams.TravelingMethod)))
	}
	if strings.TrimSpace(apiInputParams.TripNature) != "" {
		b.WriteString(fmt.Sprintf("- Trip nature: %s\n", strings.TrimSpace(apiInputParams.TripNature)))
	}
	if apiInputParams.PersonCount > 0 {
		b.WriteString(fmt.Sprintf("- Number of travelers: %d\n", apiInputParams.PersonCount))
	}
	if strings.TrimSpace(apiInputParams.GroupDemographic) != "" {
		b.WriteString(fmt.Sprintf("- Group demographic: %s\n", strings.TrimSpace(apiInputParams.GroupDemographic)))
	}
	if strings.TrimSpace(apiInputParams.Budget) != "" {
		b.WriteString(fmt.Sprintf("- Budget preference: %s\n", strings.TrimSpace(apiInputParams.Budget)))
	}
	if strings.TrimSpace(apiInputParams.CustomRequirement) != "" {
		b.WriteString(fmt.Sprintf("- Custom requirements / constraints: %s\n", strings.TrimSpace(apiInputParams.CustomRequirement)))
	}
	if strings.TrimSpace(apiInputParams.PreferredLocation) != "" {
		b.WriteString(fmt.Sprintf("- Preferred location / region: %s\n", strings.TrimSpace(apiInputParams.PreferredLocation)))
	}

	// Core instruction (focused, no response-format specifics)
	b.WriteString("\nInstructions:\n")
	b.WriteString("- Using the inputs above, suggest suitable trip locations tailored to these preferences.\n")
	b.WriteString("- Explicitly consider real-time conditions for the travel dates (current and forecast weather, official travel advisories, major local events, transport disruptions/strikes, health alerts, and natural-disaster warnings).\n")
	b.WriteString("- Flag any location with significant risk or likely major inconvenience, and give a brief reason for the flag.\n")
	b.WriteString("- Prioritize explicit user preferences (preferred_location, trip_nature, group_demographic, custom_requirement, budget) when recommending.\n")
	b.WriteString("- Provide concise rationale for each suggested location; do not add unnecessary preamble. (Output formatting will be specified separately.)\n")

	return b.String()
}

func CreateApplicationLogs(ginCtx *gin.Context, apiInputParams BindingInputParams, apiResponse ApiResponse) {
}