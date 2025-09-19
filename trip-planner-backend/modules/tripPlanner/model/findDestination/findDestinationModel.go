package findDestinationModel

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"trip-planner-backend/utilities/globalFunctions"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
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

type GeneratedDestination struct {
	Place       string `json:"place"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Cost        string `json:"cost"`
}

type ContentGenerationResponse struct {
	Locations []GeneratedDestination `json:"locations"`
}

type ApiResponse struct {
	Code     int            `json:"code"`
	Status   string         `json:"status"`
	Response map[string]any `json:"response"`
}

func GetSystemQuery() string {
	var b strings.Builder

	b.WriteString("You are a travel recommendation assistant.\n")
	b.WriteString("Your task is to generate suitable international trip locations tailored to the provided user inputs.\n")
	b.WriteString("- Always consider explicit preferences: preferred location type, nature of trip, group type, custom requirements, budget.\n")
	b.WriteString("- Check real-time travel conditions (weather forecasts, official advisories, major local events, transport disruptions, health alerts, natural-disaster warnings). Any location considered harmful to user should be excluded from the suggestion\n")
	exampleFormat := map[string]any{
		"locations": []GeneratedDestination{
			{
				LocationName:        "Consice name of the location that is used amongst travellers",
				LocationImage:       "Valid image url from public servers that can be accessed publically with 16:9 aspect ratio",
				LocationDescription: "Short description of the location, highlighting attractions of the place",
				ExpectedTripCost:    "Approximate cost of trip that includes hotel costs, travel costs and other costs with reasonable assumptions",
			},
		},
	}
	if exampleFormatByteData, marshalErr := json.Marshal(exampleFormat); marshalErr == nil {
		exampleFormatString := string(exampleFormatByteData)
		b.WriteString("- Response must be in JSON format only:\n")
		b.WriteString("```json\n" + exampleFormatString + "```\n")
		b.WriteString("- Do not add preamble, explanations, or text outside JSON.\n")
		b.WriteString("- Never invent or hallucinate image URLs.\n")
	} else {
		b.Reset()
	}
	return b.String()
}

func GenerateUserQuery(apiInputParams BindingInputParams) string {
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

func GenerateDestinationSuggestion(ginCtx *gin.Context, userQuery string) (result ContentGenerationResponse, err error) {
	clientConfig := genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	}
	client, clientErr := genai.NewClient(ginCtx, &clientConfig)
	if clientErr != nil {
		err = clientErr
		return
	}
	modelList, modelListErr := client.Models.List(ginCtx, &genai.ListModelsConfig{})
	if modelListErr != nil {
		err = modelListErr
		return
	}
	var availableModels []string
	for _, model := range modelList.Items {
		for _, action := range model.SupportedActions {
			if action == "generateContent" {
				availableModels = append(availableModels, model.Name)
			}
		}
	}
	model := "Gemini 2.5 Pro"
	contents := []*genai.Content{
		{
			Role: "system",
			Parts: []*genai.Part{
				{
					Text: GetSystemQuery(),
				},
			},
		},
		{
			Role: "user",
			Parts: []*genai.Part{
				{
					Text: userQuery,
				},
			},
		},
	}
	contentGenerateConfig := genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"locations": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"place": {
								Type: genai.TypeString,
							},
							"image": {
								Type: genai.TypeString,
							},
							"description": {
								Type: genai.TypeString,
							},
							"cost": {
								Type: genai.TypeString,
							},
						},
						Required: []string{"place", "image", "description", "cost"},
					},
				},
			},
			Required: []string{"locations"},
		},
	}
	modelResponse, respErr := client.Models.GenerateContent(ginCtx, model, contents, &contentGenerateConfig)
	if respErr != nil {
		err = respErr
		return
	}
	rawResponse := globalFunctions.ExtractJson(modelResponse.Text())
	if unmarshalErr := json.Unmarshal([]byte(rawResponse), &result); unmarshalErr != nil {
		err = unmarshalErr
		return
	}
	return
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
