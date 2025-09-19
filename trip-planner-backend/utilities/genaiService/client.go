package genaiService

import (
	"context"
	"os"

	"google.golang.org/genai"
)

var genaiClient *genai.Client

func GetClient() (client *genai.Client, err error) {
	if genaiClient != nil {
		client = genaiClient
		return
	}
	clientConfig := genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	}
	client, clientErr := genai.NewClient(context.Background(), &clientConfig)
	if clientErr != nil {
		client = nil
		err = clientErr
		return
	}
	genaiClient = client
	return
}
