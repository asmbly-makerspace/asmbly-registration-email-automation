package mailservice

import (
	"errors"

	"github.com/mailjet/mailjet-apiv3-go/v3/resources"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

var ErrNotFound = errors.New("couldn't find any templates with that name")
var ErrMultipleFound = errors.New("found more than one template with that name")

const EmailSubject = "Information about your upcoming class at Asmbly"
const FromEmail = "membership@asmbly.org"
const FromName = "Asmbly Education Team"

type MJCredentials struct {
	PublicKey, SecretKey string
}

type mailServiceInterface interface {
	SendEmail(values EmailInfo) error
	GetTemplateIDByName(name string) (int, error)
}

type Client struct {
	mailService mailServiceInterface
}

func NewClient(mailService mailServiceInterface) *Client {
	return &Client{mailService}
}

func (c *Client) SendRegistrationEmail(className, email, firstName string) error {

	templateID, err := c.mailService.GetTemplateIDByName(className)
	if err != nil {
		return err
	}

	values := EmailInfo{
		Subject:    EmailSubject,
		From:       FromEmail,
		FromName:   FromName,
		To:         []string{email},
		TemplateID: templateID,
		Variables: map[string]interface{}{
			"firstname":  firstName,
			"class_name": className,
		},
	}
	err = c.mailService.SendEmail(values)
	if err != nil {
		return err
	}
	return nil
}

type EmailInfo struct {
	Subject    string
	From       string
	FromName   string
	To         []string
	TemplateID int
	Variables  map[string]interface{}
}

type MJMailService struct {
	client *mailjet.Client
}

func NewMJClient(creds MJCredentials) *MJMailService {
	client := mailjet.NewMailjetClient(creds.PublicKey, creds.SecretKey)
	return &MJMailService{
		client: client,
	}
}

func (mj *MJMailService) SendEmail(values EmailInfo) error {

	var recipients mailjet.RecipientsV31

	for _, recipient := range values.To {
		recipients = append(recipients, mailjet.RecipientV31{Email: recipient})
	}

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: values.From,
				Name:  values.FromName,
			},
			To:               &recipients,
			TemplateID:       values.TemplateID,
			TemplateLanguage: true,
			Subject:          values.Subject,
			Variables:        values.Variables,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	_, err := mj.client.SendMailV31(&messages)
	if err != nil {
		return err
	}
	return nil
}

func (mj *MJMailService) GetTemplateIDByName(name string) (id int, err error) {

	var data []resources.Template

	nameFilter := mailjet.Filter("Name", name)
	count, _, err := mj.client.List("template", &data, nameFilter)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, ErrNotFound
	} else if count > 1 {
		return 0, ErrMultipleFound
	}

	return int(data[0].ID), nil
}
