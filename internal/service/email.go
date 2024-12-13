package service

import (
	"errors"
	"meta-node-ficam/internal/request"
	"meta-node-ficam/internal/repositories"
	"meta-node-ficam/utils"
	"time"
)


type EmailService interface {
	EmailVerification(request request.EmailVerificationRequest) error
	EmailAuthentication(request request.EmailAuthenticationRequest) error
}

type emailService struct {
	emailRepo repositories.EmailRepository
}

func NewEmailService(
	emailRepo repositories.EmailRepository,
) EmailService {
	return &emailService{emailRepo}
}

func (svc *emailService) EmailVerification(request request.EmailVerificationRequest) error {
	codeEmail := utils.GenerateVerificationToken()
	exipredTime := time.Now().Add(time.Duration(utils.VerificationExpiredTime) * time.Minute)
	err := svc.emailRepo.SaveEmailVerification(request.Email,codeEmail, exipredTime)
	if err != nil {
		return err
	}
	go func() {
		utils.ReplyEmailRegister(request.Email, "Your code to verify email is : " + codeEmail,"Email Verifcation Metanode")
	}()
	return nil
}

func (svc *emailService) EmailAuthentication(request request.EmailAuthenticationRequest) error {
	existsEmail, err := svc.emailRepo.CheckEmailVerification(request.Email,request.CodeEmail, time.Now())
	if err != nil {
		return err
	}
	if !existsEmail {
		return errors.New("Invalid email, code or code expired")
	}
	
	go func() {
		utils.ReplyEmailRegister(request.Email, "Your email has been successfully authenticated","Email Verification Metanode")
	}()

	return nil
}
