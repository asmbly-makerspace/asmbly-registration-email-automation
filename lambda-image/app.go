package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mkmiller6/asmbly-class-email-automation/mailservice"
)

func lambdaHandler(event Event) error {
	// Create the mail service
	creds := mailservice.MJCredentials{
		PublicKey: os.Getenv("MJ_APIKEY_PUBLIC"),
		SecretKey: os.Getenv("MJ_APIKEY_PRIVATE"),
	}

	mjClient := mailservice.NewMJClient(creds)
	mailService := mailservice.NewClient(mjClient)

	var eventJson EventBody
	json.Unmarshal([]byte(event.Body), &eventJson)

	className := neonClient.GetClassByID(eventJson.Data.EventID)
	registrantAcct := neonClient.GetAccountByID(eventJson.Data.RegistrantAccountID)

	err := mailService.SendRegistrationEmail(
		className,
		registrantAcct.IndividualAccount.PrimaryContact.Email1,
		eventJson.Data.Tickets[0].Attendees[0].FirstName,
	)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	lambda.Start(lambdaHandler)
}
