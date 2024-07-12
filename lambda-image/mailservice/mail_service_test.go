package mailservice_test

import "testing"

func TestSendRegistrationEmail(t *testing.T) {
	mailService := MockMailService{}

	className := "Test Filament 3D Printing"
	email := "testing@example.com"

	err := mailService.SendRegistrationEmail(className, email)

	if err != nil {
		t.Fatalf("Sending registraiton email failed with error %s", err)
	}
}

type MockMailService struct{}

func (m *MockMailService) SendRegistrationEmail(className, email string) error {
	return nil
}
