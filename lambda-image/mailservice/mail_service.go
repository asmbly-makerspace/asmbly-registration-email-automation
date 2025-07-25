package mailservice

import (
	"errors"
	"log"

	"github.com/mailjet/mailjet-apiv3-go/v3/resources"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

var ErrNotFound = errors.New("couldn't find any templates with that name")
var ErrMultipleFound = errors.New("found more than one template with that name")

const EmailSubject = "Information about your upcoming class at Asmbly"
const FromEmail = "classes@asmbly.org"
const FromName = "Asmbly Education Team"

type MJCredentials struct {
	PublicKey, SecretKey string
}

type mailServiceInterface interface {
	SendEmail(values EmailInfo) error
	GetTemplateIDByName(name string) (count int, id int, err error)
}

type Client struct {
	mailService mailServiceInterface
}

func NewClient(mailService mailServiceInterface) *Client {
	return &Client{mailService}
}

func (c *Client) SendRegistrationEmail(className, email, firstName string) error {

	count, templateID, err := c.mailService.GetTemplateIDByName(className)
	if err != nil {
		return err
	} else if count == 0 {
		log.Printf("no email template found for class %q", className)
		return nil
	}

	values := EmailInfo{
		Subject:    EmailSubject,
		From:       FromEmail,
		FromName:   FromName,
		To:         []string{email},
		TemplateID: templateID,
		Variables: map[string]interface{}{
			"first_name": firstName,
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

func (mj *MJMailService) GetTemplateIDByName(name string) (count, id int, err error) {

	var data []resources.Template

	nameFilter := mailjet.Filter("Name", name)
	count, _, err = mj.client.List("template", &data, nameFilter)
	if err != nil {
		return 0, 0, err
	}
	if count == 0 {
		return count, 0, nil
	} else if count > 1 {
		return count, 0, ErrMultipleFound
	}

	return count, int(data[0].ID), nil
}
