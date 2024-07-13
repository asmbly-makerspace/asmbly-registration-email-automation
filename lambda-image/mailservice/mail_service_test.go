package mailservice_test

import (
	"testing"

	"github.com/mkmiller6/asmbly-class-email-automation/mailservice"
)

func TestSendRegistrationEmail(t *testing.T) {
	t.Run("one template with that name in db", func(t *testing.T) {
		fakeService := &fakeMailService{
			templateDatabase: []fakeTemplateEntry{
				{ID: 123, Name: "Test Filament 3D Printing"},
			},
		}
		mailClient := mailservice.NewClient(fakeService)

		className := "Test Filament 3D Printing"
		email := "testing@example.com"
		firstName := "John"

		err := mailClient.SendRegistrationEmail(className, email, firstName)

		if err != nil {
			t.Fatalf("Sending registraiton email failed with error %s", err)
		}

		got := fakeService.calls
		want := 1

		if got != want {
			t.Errorf("got %d send email calls, but wanted %d", got, want)
		}

		templateGot := fakeService.calledWith.TemplateID
		templateWant := 123

		if templateGot != templateWant {
			t.Errorf("got templateId %d, want %d", templateGot, templateWant)
		}
	})
	t.Run("multiple templates with that name in db", func(t *testing.T) {
		fakeService := &fakeMailService{
			templateDatabase: []fakeTemplateEntry{
				{ID: 123, Name: "Test Filament 3D Printing"},
				{ID: 456, Name: "Test Filament 3D Printing"},
			},
		}
		mailClient := mailservice.NewClient(fakeService)

		className := "Test Filament 3D Printing"
		email := "testing@example.com"
		firstName := "John"

		err := mailClient.SendRegistrationEmail(className, email, firstName)

		if err == nil {
			t.Fatal("expected an error but didn't get one")
		}
		assertError(t, err, mailservice.ErrMultipleFound)

		got := fakeService.calls
		want := 0

		if got != want {
			t.Errorf("got %d send email calls, but wanted %d", got, want)
		}
	})
	t.Run("no templates with name in db", func(t *testing.T) {
		fakeService := &fakeMailService{
			templateDatabase: []fakeTemplateEntry{},
		}
		mailClient := mailservice.NewClient(fakeService)

		className := "Test Filament 3D Printing"
		email := "testing@example.com"
		firstName := "John"

		err := mailClient.SendRegistrationEmail(className, email, firstName)

		if err == nil {
			t.Fatal("expected an error but didn't get one")
		}
		assertError(t, err, mailservice.ErrNotFound)

		got := fakeService.calls
		want := 0

		if got != want {
			t.Errorf("got %d send email calls, but wanted %d", got, want)
		}
	})
}

type fakeMailService struct {
	calls            int
	templateDatabase []fakeTemplateEntry
	calledWith       mailservice.EmailInfo
}

func (f *fakeMailService) SendEmail(values mailservice.EmailInfo) error {
	f.calls++
	f.calledWith = values
	return nil
}

func (f *fakeMailService) GetTemplateIDByName(name string) (id int, err error) {
	var found []fakeTemplateEntry
	for _, template := range f.templateDatabase {
		if template.Name == name {
			found = append(found, template)
		}
	}

	if len(found) == 0 {
		return 0, mailservice.ErrNotFound
	} else if len(found) > 1 {
		return 0, mailservice.ErrMultipleFound
	}

	return int(found[0].ID), nil
}

type fakeTemplateEntry struct {
	ID   int64
	Name string
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %q, want %q", got, want)
	}
}
