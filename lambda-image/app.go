package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mkmiller6/asmbly-class-email-automation/mailservice"
)

func lambdaHandler(ctx context.Context) error {
	// Create the mail service
	publicKey := os.Getenv("MJ_APIKEY_PUBLIC")
	secretKey := os.Getenv("MJ_APIKEY_PRIVATE")

	creds := mailservice.Credentials{
		PublicKey: publicKey,
		SecretKey: secretKey,
	}

	mailService := mailservice.NewClient(creds)

	mailService.SendRegistrationEmail("", "")

	return nil
}

func main() {
	lambda.Start(lambdaHandler)
}
