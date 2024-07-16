package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mkmiller6/asmbly-class-email-automation/mailservice"
	"github.com/mkmiller6/neon-go-client"
	"github.com/mkmiller6/neon-go-client/client"
)

type AccountsInterface interface {
	GetByID(int) (*neon.Account, error)
}

type EventsInterface interface {
	GetByID(int) (*neon.Event, error)
}

type NeonResult struct {
	Response string
	Error    error
}

func GetNeonEventName(client EventsInterface, eventJson *EventBody, c chan<- NeonResult) {

	eventId, err := strconv.Atoi(eventJson.Data.EventID)
	if err != nil {
		c <- NeonResult{
			Response: "",
			Error:    fmt.Errorf("parsing neon event id string failed with error: %q", err),
		}
		return
	}
	neonEvent, err := client.GetByID(eventId)
	if err != nil {
		c <- NeonResult{
			Response: "",
			Error:    fmt.Errorf("getting neon event failed with error: %q", err),
		}
		return
	}

	c <- NeonResult{
		Response: strings.Split(neonEvent.Name, " w/")[0],
		Error:    nil,
	}
}

func GetNeonAccountEmail(client AccountsInterface, eventJson *EventBody, c chan<- NeonResult) {
	neonAcctId, err := strconv.Atoi(eventJson.Data.RegistrantAccountID)
	if err != nil {
		c <- NeonResult{
			Error:    fmt.Errorf("parsing neon account id string failed with error: %q", err),
			Response: "",
		}
		return
	}
	registrantAcct, err := client.GetByID(neonAcctId)
	if err != nil {
		c <- NeonResult{
			Error:    fmt.Errorf("getting neon account failed with error: %q", err),
			Response: "",
		}
		return
	}

	c <- NeonResult{
		Response: registrantAcct.IndividualAccount.PrimaryContact.Email1,
		Error:    nil,
	}
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

	eventChan := make(chan NeonResult)
	acctChan := make(chan NeonResult)

	go GetNeonEventName(neonClient.Events, &eventJson, eventChan)
	go GetNeonAccountEmail(neonClient.Accounts, &eventJson, acctChan)

	eventResult := <-eventChan
	acctResult := <-acctChan

	if eventResult.Error != nil {
		log.Fatalf("neon query failed with error: %q", eventResult.Error)
	} else if acctResult.Error != nil {
		log.Fatalf("neon query failed with error: %q", acctResult.Error)
	}

	neonEventName := eventResult.Response
	registrantEmail := acctResult.Response

	if os.Getenv("DEV") == "1" && registrantEmail != os.Getenv("TEST_EMAIL") {
		return nil
	}

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
