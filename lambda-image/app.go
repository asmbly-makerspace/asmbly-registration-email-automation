package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mkmiller6/asmbly-class-email-automation/mailservice"
	"github.com/mkmiller6/neon-go-client"
	"github.com/mkmiller6/neon-go-client/client"
)

func getNeonEventName(client *client.API, eventJson *EventBody) (name string) {

	eventId, err := strconv.Atoi(eventJson.Data.EventID)
	if err != nil {
		log.Fatalf("parsing neon event id string failed with error: %q", err)
	}
	neonEvent, err := client.Events.GetByID(eventId)
	if err != nil {
		log.Fatalf("getting neon event failed with error: %q", err)
	}

	return strings.Split(neonEvent.Name, " w/")[0]
}

func getNeonAccountEmail(client *client.API, eventJson *EventBody) (email string) {
	neonAcctId, err := strconv.Atoi(eventJson.Data.RegistrantAccountID)
	if err != nil {
		log.Fatalf("parsing neon account id string failed with error: %q", err)
	}
	registrantAcct, err := client.Accounts.GetByID(neonAcctId)
	if err != nil {
		log.Fatalf("getting neon account failed with error: %q", err)
	}

	return registrantAcct.IndividualAccount.PrimaryContact.Email1
}

func lambdaHandler(event Event) error {
	// Create the mail service
	mjCreds := mailservice.MJCredentials{
		PublicKey: os.Getenv("MJ_APIKEY_PUBLIC"),
		SecretKey: os.Getenv("MJ_APIKEY_PRIVATE"),
	}

	mjClient := mailservice.NewMJClient(mjCreds)
	mailService := mailservice.NewClient(mjClient)

	// Create the Neon client
	neonBackend := neon.GetBackendWithConfig()
	neonClient := client.New(os.Getenv("NEON_APIUSER"), os.Getenv("NEON_APIKEY"), neonBackend)

	var eventJson EventBody
	json.Unmarshal([]byte(event.Body), &eventJson)

	neonEventName := getNeonEventName(neonClient, &eventJson)
	registrantEmail := getNeonAccountEmail(neonClient, &eventJson)

	err := mailService.SendRegistrationEmail(
		neonEventName,
		registrantEmail,
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
