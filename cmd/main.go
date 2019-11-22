package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/akillmer/cloudmail"
)

func main() {
	// this key is explicity for testing as provided by Google
	os.Setenv("RECAPTCHA_SECRET", "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe")
	funcframework.RegisterHTTPFunction("/", cloudmail.SendMessage)
	port := "9000"

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
