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
	"gopkg.in/gomail.v2"
)

func ReplyEmailRegister(emailReceiver, bodyMessage , subject string) error {
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatal("Failed to load config in Reply email: ", err)
	}
	emailSender := model.Sender{
		Address : cfg.AWSConfig.SenderAddress,
		Subject : subject,
	}
	destinations := []*string{aws.String(emailReceiver)}
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSConfig.Region),
		Credentials: credentials.NewStaticCredentials(
			cfg.AWSConfig.CredentialsConfig.Id,
			cfg.AWSConfig.CredentialsConfig.SecretKey,
			cfg.AWSConfig.CredentialsConfig.Token),
	})
	svc := ses.New(sess)
	rawEmail := &ses.RawMessage{
		Data: formatContentEmailToSend(emailSender,emailReceiver, bodyMessage),
	}
	_, err = svc.SendRawEmail(&ses.SendRawEmailInput{
		Source:       aws.String(emailSender.Address),
		Destinations: destinations,
		RawMessage:   rawEmail,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func formatContentEmailToSend(emailSender model.Sender,emailReceiver, bodyMessage string) []byte {
	m := gomail.NewMessage()
	m.SetHeader("From", emailSender.Address)
	m.SetHeader("Subject", emailSender.Subject)
	m.SetHeader("To", emailReceiver)
	m.SetBody("text/html", bodyMessage)

	var buffer bytes.Buffer
	_, err := m.WriteTo(&buffer)
	if err != nil {
		println(err.Error())
	}
	return buffer.Bytes()
}
