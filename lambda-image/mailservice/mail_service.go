package mailservice

// Import the Mailjet wrapper
import (
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type MailService interface {
	SendRegistrationEmail(className, email string) error
}

func NewClient() *mailjet.Client {
	// Get your environment Mailjet keys and connect
	publicKey := os.Getenv("MJ_APIKEY_PUBLIC")
	secretKey := os.Getenv("MJ_APIKEY_PRIVATE")

	mj := mailjet.NewMailjetClient(publicKey, secretKey)

	return mj
}
