package generateItineraryModel

import (
	"encoding/json"
	"fmt"
	"strings"
	globalconstant "trip-planner-backend/globalConstant"
	"trip-planner-backend/utilities/genaiService"
	"trip-planner-backend/utilities/globalFunctions"

	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

type ApiInputParams struct {
	UserId           int    `json:"user_id" binding:"required"`
	UserLocation     string `json:"user_location" binding:"required"`
	Destination      string `json:"destination" binding:"required"`
	TravelDays       int    `json:"travel_days"`
	TravelDateTime   string `json:"travel_date_time"`
	PersonCount      int    `json:"person_count"`
	GroupDemographic string `json:"group_demographic"`
}

type ApiResponse struct {
	Code     int                       `json:"code"`
	Status   string                    `json:"status"`
	Error    string                    `json:"error"`
	Response ContentGenerationResponse `json:"response"`
}

type ContentGenerationResponse struct {
	Itineraries []Itinerary `json:"itineraries"`
	Summary     string      `json:"summary"`
}

type Itinerary struct {
	Overview     string       `json:"overview"`
	StartDate    string       `json:"start_date"`
	EndDate      string       `json:"end_date"`
	Transport    []Transport  `json:"transport"`
	Accomodation Accomodation `json:"accomodation"`
	Guide        bool         `json:"guide"`
	Photoshoot   bool         `json:"photoshoot"`
}

type Transport struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Distance    string `json:"distance"`
}

type Accomodation struct {
	Duration int `json:"duration"`
	Area string `json:"area"`
}

func GetSystemQuery() string {
	var b strings.Builder

	b.WriteString("You are a travel itinerary generator expert.\n")
	b.WriteString("Your task is to generate suitable itinerary based on user selected destination and other constraints such as number of days, start date, person count and demographic of the group .\n")
	b.WriteString("Generate different itinerary parts based on logical seperation of the trip and its required time and travel\n")
	b.WriteString("Check real-time travel conditions (weather forecasts, official advisories, major local events, transport disruptions, health alerts, natural-disaster warnings). Any location considered harmful to user should be excluded from the suggestion\n")
	b.WriteString("Strictly adhere to the constraints provided by the user. Any part of trip should not overlap or there should be no time gaps between different parts of the trip and there should be no unreasonable travel distance between any parts of the trip.\n")
	b.WriteString("Trip should also include travel from user location to destination and should also consider return trip to the user location.\n")
	exampleFormat := &ContentGenerationResponse{
		Itineraries: []Itinerary{
			{
				Overview:  "General Overview of this part of the trip that explains what to explore here",
				StartDate: "Start date of this part of trip in dd-mm-yy format",
				EndDate:   "End date of this part of trip in dd-mm-yy format",
				Transport: []Transport{
					{
						Origin:      "Origin place where user was before starting this part of the trip",
						Destination: "Location for this part of the trip",
						Distance:    "Distance between origin and destination in kms",
					},
				},
				Accomodation: Accomodation{
					Duration: 3,
					Area: "Local area for the current part of the trip where we might need to book hotel",
				},
				Guide:      true,
				Photoshoot: true,
			},
		},
		Summary: "General summary of the whole trip and list of things to do during the trip that might attract the users",
	}
	exampleFormatBytes, _ := json.Marshal(exampleFormat)
	exampleFormatString := string(exampleFormatBytes)
	b.WriteString("- Response should be strictly in below JSON format: ```json\n" + exampleFormatString + "```\n")
	b.WriteString("- Location name should be consice and must be most commnly used name for that place.\n")
	b.WriteString("- Accomodation should consider on days when user may require to stay in hotel of that part of the trip and it should keep the area of accomodation consistent around the tourism places\n")
	b.WriteString("- It should suggest if that part of the trip should suggest guide or photoshoot only if that is possible\n")
	b.WriteString("- Summary should match with the different part of the trip and should highlight attractions of the trip\n")
	b.WriteString("- Overview should also suggest some suggestions of what to do or explore or eat at that location")
	return b.String()
}

func GenerateUserQuery(apiInputParams ApiInputParams) string {
	var b strings.Builder

	b.WriteString("Given the user inputs below, suggest suitable trip divided into logical parts tailored to user preferences and constraints.\n\n")
	b.WriteString("User inputs:\n")
	b.WriteString(fmt.Sprintf("- User Start Location: %s\n", strings.TrimSpace(apiInputParams.UserLocation)))
	b.WriteString(fmt.Sprintf("- Trip Destination: %s\n", strings.TrimSpace(apiInputParams.Destination)))

	// Add other fields only when they are non-zero / non-empty
	if apiInputParams.TravelDays > 0 {
		b.WriteString(fmt.Sprintf("- Travel Days: %d\n", apiInputParams.TravelDays))
	}
	if strings.TrimSpace(apiInputParams.TravelDateTime) != "" {
		b.WriteString(fmt.Sprintf("- Travel Start Date: %s\n", strings.TrimSpace(apiInputParams.TravelDateTime)))
	}
	if apiInputParams.PersonCount > 0 {
		b.WriteString(fmt.Sprintf("- Person Count: %d\n", apiInputParams.PersonCount))
	}
	if strings.TrimSpace(apiInputParams.GroupDemographic) != "" {
		b.WriteString(fmt.Sprintf("- Group Demographic: %s\n", strings.TrimSpace(apiInputParams.GroupDemographic)))
	}

	// Core instruction (focused, no response-format specifics)
	b.WriteString("\nInstructions:\n")
	b.WriteString("- Using the inputs above, suggest single cohesive itinerary consisting of different parts.\n")
	b.WriteString("- Explicitly consider real-time conditions for the travel dates (current and forecast weather, official travel advisories, major local events, transport disruptions/strikes, health alerts, and natural-disaster warnings).\n")
	b.WriteString("- Exclude any location with significant risk or likely major inconvenience, and give a brief warning in overview of that part of the trip.\n")
	b.WriteString("- Prioritize explicit user preferences (travel days, travel start date, group_demographic) when recommending.\n")
	b.WriteString("- Provide concise trip itinerary with proper rationale for each suggested part of the trip; do not add unnecessary preamble.\n")

	return b.String()
}

func GenerateItinerarySuggestion(ginCtx *gin.Context, userQuery string) (result ContentGenerationResponse, err error) {
	client, clientErr := genaiService.GetClient()
	if clientErr != nil {
		err = clientErr
		return
	}

	model := globalconstant.GEMINI_MODEL
	if isModelAvailableErr := genaiService.IsModelAvailable(ginCtx, client, model); isModelAvailableErr != nil {
		err = isModelAvailableErr
		return
	}
	contents := []*genai.Content{
		{
			Role: "model",
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
		ResponseMIMEType: "text/plain",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"itineraries": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"overview": {
								Type: genai.TypeString,
							},
							"start_date": {
								Type: genai.TypeString,
							},
							"end_date": {
								Type: genai.TypeString,
							},
							"transport": {
								Type: genai.TypeArray,
								Items: &genai.Schema{
									Type: genai.TypeObject,
									Properties: map[string]*genai.Schema{
										"origin": {
											Type: genai.TypeString,
										},
										"destination": {
											Type: genai.TypeString,
										},
										"distance": {
											Type: genai.TypeString,
										},
									},
								},
							},
							"accomodation": {
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"duration": {
										Type: genai.TypeInteger,
									},
									"area": {
										Type: genai.TypeString,
									},
								},
							},
							"guide": {
								Type: genai.TypeBoolean,
							},
							"photoshoot": {
								Type: genai.TypeBoolean,
							},
						},
						Required: []string{"overview", "start_date", "end_date", "transport", "accomodation", "guide", "photoshoot"},
					},
				},
				"summary": {
					Type: genai.TypeString,
				},
			},
			Required: []string{"itineraries", "summary"},
		},
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		},
	}
	modelResponse, respErr := client.Models.GenerateContent(ginCtx, model, contents, &contentGenerateConfig)
	if respErr != nil {
		err = respErr
		return
	}
	rawResponse, _ := globalFunctions.ExtractJson(modelResponse.Text())
	if unmarshalErr := json.Unmarshal([]byte(rawResponse), &result); unmarshalErr != nil {
		err = unmarshalErr
		return
	}
	return
}

func CreateApplicationLogs(ginCtx *gin.Context, apiInputParams ApiInputParams, apiResponse ApiResponse) {

	fileName := "generate_itinerary"

	logData := map[string]any{}
	logData["user_id"] = globalFunctions.ConvertValueToString(apiInputParams.UserId)
	logData["user_location"] = globalFunctions.ConvertValueToString(apiInputParams.UserLocation)
	logData["destination"] = globalFunctions.ConvertValueToString(apiInputParams.Destination)
	logData["travel_days"] = globalFunctions.ConvertValueToString(apiInputParams.TravelDays)
	logData["travel_date_time"] = globalFunctions.ConvertValueToString(apiInputParams.TravelDateTime)
	logData["person_count"] = globalFunctions.ConvertValueToString(apiInputParams.PersonCount)
	logData["group_demographic"] = globalFunctions.ConvertValueToString(apiInputParams.GroupDemographic)

	logData["code"] = globalFunctions.ConvertValueToString(apiResponse.Code)
	logData["status"] = globalFunctions.ConvertValueToString(apiResponse.Status)
	logData["response"] = globalFunctions.ConvertJsonValToString(apiResponse.Response)

	globalFunctions.WriteJsonLogs(ginCtx, fileName, logData)
}