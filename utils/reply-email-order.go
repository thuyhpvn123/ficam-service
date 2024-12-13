package utils

import (
	"bytes"
	"fmt"
	"log"
	"meta-node-ficam/internal/config"
	"meta-node-ficam/internal/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	// "gopkg.in/gomail.v2"
	"html/template"

)


func ReplyEmailOrder(emailReceiver string, data model.Data , subject string) error {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatal("Failed to load config in Reply email: ", err)
	}
	emailSender := model.Sender{
		Address : cfg.AWSConfig.SenderAddress,
		Subject : subject,
	}
	// destinations := []*string{aws.String(emailReceiver)}
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSConfig.Region),
		Credentials: credentials.NewStaticCredentials(
			cfg.AWSConfig.CredentialsConfig.Id,
			cfg.AWSConfig.CredentialsConfig.SecretKey,
			cfg.AWSConfig.CredentialsConfig.Token),
	})
	svc := ses.New(sess)
	recipient := model.Recipient{
		ToEmails:  []string{emailReceiver},
		CcEmails:  []string{},
		BccEmails: []string{cfg.EmailAdmin},
	}
	var recipients []*string
	for _, r := range recipient.ToEmails {
		recipient := r
		recipients = append(recipients, &recipient)
	}
	var ccRecipients []*string
	if len(recipient.CcEmails) > 0 {
		for _, r := range recipient.CcEmails {
			ccrecipient := r
			ccRecipients = append(ccRecipients, &ccrecipient)
		}
	}
	var bccRecipients []*string
	if len(recipient.BccEmails) > 0 {
		for _, r := range recipient.BccEmails {
			bccrecipient := r
			recipients = append(recipients, &bccrecipient)
		}
	}
	htmlBody, err := parseHTMLTemplate(cfg.TemplateEmailOrderPath, data)
	if err != nil {
		log.Fatalf("unable to parse template, %v", err)
		return err
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses:  ccRecipients,
			ToAddresses:  recipients,
			BccAddresses: bccRecipients,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(htmlBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(emailSender.Address),
	}
	// Send the email
	_, err = svc.SendEmail(input)
	if err != nil {
		log.Fatalf("failed to send email, %v", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}


// parseHTMLTemplate parses the HTML template and returns it as a string
func parseHTMLTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}

	var result string
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return "", err
	}

	result = buf.String()
	return result, nil
}
