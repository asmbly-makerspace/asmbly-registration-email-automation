package main_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	app "github.com/mkmiller6/asmbly-class-email-automation"
	"github.com/mkmiller6/neon-go-client"
	"github.com/stretchr/testify/assert"
)

func TestGetEventName(t *testing.T) {
	client := fakeEventsClient{
		EventsDB: []neonEvent{
			{
				Id:   7008,
				Name: "Test Event w/ Test Teacher",
			},
			{
				Id:   8001,
				Name: "Test Event 2 w/ Another Teacher",
			},
		},
	}
	t.Run("event in db", func(t *testing.T) {
		json := &app.EventBody{
			Data: app.EventData{
				EventID: "7008",
			},
		}

		c := make(chan app.NeonResult)

		go app.GetNeonEventName(client, json, c)
		result := <-c
		if result.Error != nil {
			t.Fatalf("error finding event in database: %q", result.Error)
		}

		got := result.Response
		want := "Test Event"

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
	t.Run("event not in db", func(t *testing.T) {
		json := &app.EventBody{
			Data: app.EventData{
				EventID: "7010",
			},
		}

		c := make(chan app.NeonResult)

		go app.GetNeonEventName(client, json, c)
		result := <-c
		if result.Error == nil {
			t.Fatal("expected an error but didn't get one")
		}
		expectedError := fmt.Sprintf("getting neon event failed with error: %q", ErrEventNotFound)
		assert.EqualErrorf(t, result.Error, expectedError, "Error should be: %v, got: %v", expectedError, result.Error)

		got := result.Response
		want := ""

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
func TestGetAccountEmail(t *testing.T) {
	client := fakeAccountsClient{
		AccountsDB: []neonAccount{
			{
				Id: 89,
				IndividualAccount: neon.IndividualAccount{
					PrimaryContact: &neon.Contact{
						Email1: "testemail@example.com",
					},
				},
			},
			{
				Id: 1500,
				IndividualAccount: neon.IndividualAccount{
					PrimaryContact: &neon.Contact{
						Email1: "testemail2@example.com",
					},
				},
			},
		},
	}
	t.Run("acct in db", func(t *testing.T) {
		json := &app.EventBody{
			Data: app.EventData{
				RegistrantAccountID: "89",
			},
		}

		c := make(chan app.NeonResult)

		go app.GetNeonAccountEmail(client, json, c)
		result := <-c
		if result.Error != nil {
			t.Fatalf("error finding acct in database: %q", result.Error)
		}

		got := result.Response
		want := "testemail@example.com"

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
	t.Run("acct not in db", func(t *testing.T) {
		json := &app.EventBody{
			Data: app.EventData{
				RegistrantAccountID: "7010",
			},
		}

		c := make(chan app.NeonResult)

		go app.GetNeonAccountEmail(client, json, c)
		result := <-c
		if result.Error == nil {
			t.Fatal("expected an error but didn't get one")
		}
		expectedError := fmt.Sprintf("getting neon account failed with error: %q", ErrAcctNotFound)
		assert.EqualErrorf(t, result.Error, expectedError, "Error should be: %v, got: %v", expectedError, result.Error)

		got := result.Response
		want := ""

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}

var ErrEventNotFound = errors.New("could not find class in database")
var ErrAcctNotFound = errors.New("could not find account in database")

type neonEvent struct {
	Id   int
	Name string
}

type fakeEventsClient struct {
	EventsDB []neonEvent
}

func (f fakeEventsClient) GetByID(id int) (*neon.Event, error) {
	for _, event := range f.EventsDB {
		if event.Id == id {
			stringId := strconv.Itoa(event.Id)
			return &neon.Event{
				Id:   stringId,
				Name: event.Name,
			}, nil
		}
	}
	return nil, ErrEventNotFound
}

type neonAccount struct {
	Id                int
	IndividualAccount neon.IndividualAccount
}

type fakeAccountsClient struct {
	AccountsDB []neonAccount
}

func (f fakeAccountsClient) GetByID(id int) (*neon.Account, error) {
	for _, acct := range f.AccountsDB {
		if acct.Id == id {
			return &neon.Account{
				IndividualAccount: &acct.IndividualAccount,
			}, nil
		}
	}
	return nil, ErrAcctNotFound
}
