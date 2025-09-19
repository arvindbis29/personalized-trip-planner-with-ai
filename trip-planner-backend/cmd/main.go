package main

import (
	"fmt"
	"net/http"
	"os"
	"trip-planner-backend/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvErr := godotenv.Load()
	if loadEnvErr != nil {
		fmt.Println("Some error in loading env file: " + loadEnvErr.Error())
	}
	port := os.Getenv("LISTEN_PORT")
	ginEngine := gin.New()
	handler.RouteRequests(ginEngine)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: ginEngine,
	}

	if serverErr := server.ListenAndServe(); serverErr != nil {
		fmt.Println("Some issue occured while initiating the server: " + serverErr.Error())
	}
}
