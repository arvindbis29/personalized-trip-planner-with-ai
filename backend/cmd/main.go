package main

import (
	"fmt"
	"net/http"
	"trip-planner-backend/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	port := "8080"
	ginEngine := gin.New()
	handler.RouteRequests(ginEngine)
	server := &http.Server{
		Addr:    ":"+port,
		Handler: ginEngine,
	}

	if serverErr := server.ListenAndServe(); serverErr != nil {
		fmt.Println("BC AAAAAAAAAGGGGG: " + serverErr.Error())
	}
}