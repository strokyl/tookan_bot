package main

import (
	"github.com/strokyl/tookan_bot/pkg/googlecloudfunctions"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"os"
	"log"
)

func main() {
	funcframework.RegisterHTTPFunction("/ChangeHourForm", googlecloudfunctions.ChangeHourForm)
	funcframework.RegisterHTTPFunction("/ChangeHour", googlecloudfunctions.ChangeHour)
	funcframework.RegisterHTTPFunction("/CancelTask", googlecloudfunctions.CancelTask)
	funcframework.RegisterHTTPFunction("/SendChangeHourSMS", googlecloudfunctions.SendChangeHourSMS)
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
