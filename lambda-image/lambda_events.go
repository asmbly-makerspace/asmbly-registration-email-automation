package main

import "time"

type Event struct {
	Body    string       `json:"body,omitempty"`
	Headers EventHeaders `json:"headers,omitempty"`
}

type EventHeaders struct {
	UserAgent     string `json:"user-agent,omitempty"`
	Authorization string `json:"authorization,omitempty"`
}

type EventBody struct {
	EventTrigger   string    `json:"eventTrigger,omitempty"`
	EventTimestamp time.Time `json:"eventTimestamp,omitempty"`
	Data           EventData `json:"data,omitempty"`
}

type EventData struct {
	EventID             string   `json:"eventId,omitempty"`
	RegistrantAccountID string   `json:"registrantAccountId,omitempty"`
	Tickets             []Ticket `json:"tickets,omitempty"`
}

type Ticket struct {
	Attendees []Attendee `json:"attendees,omitempty"`
}

type Attendee struct {
	RegistrationStatus  string `json:"registrationStatus,omitempty"`
	RegistrantAccountID string `json:"registrantAccountId,omitempty"`
	FirstName           string `json:"firstName,omitempty"`
	LastName            string `json:"lastName,omitempty"`
}
