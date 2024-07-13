package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mkmiller6/asmbly-class-email-automation/mailservice"
)

func lambdaHandler(ctx context.Context) error {
	// Create the mail service
	creds := mailservice.MJCredentials{
		PublicKey: os.Getenv("MJ_APIKEY_PUBLIC"),
		SecretKey: os.Getenv("MJ_APIKEY_PRIVATE"),
	}

	mjClient := mailservice.NewMJClient(creds)

	mailService := mailservice.NewClient(mjClient)

	mailService.SendRegistrationEmail("", "", "")

	return nil
}

func main() {
	lambda.Start(lambdaHandler)
}
