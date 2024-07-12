package mailservice

// Import the Mailjet wrapper
import (
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type Credentials struct {
	PublicKey string
	SecretKey string
}

type MailServiceInterface interface {
	SendRegistrationEmail(className, email string) error
}

func NewClient(creds Credentials) *MailService {

	mj := mailjet.NewMailjetClient(creds.PublicKey, creds.SecretKey)

	mailService := &MailService{
		client: mj,
	}

	return mailService
}

type MailService struct {
	client *mailjet.Client
}

func (ms *MailService) SendRegistrationEmail(className, email string) error {
	return nil
}
