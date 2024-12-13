package model

import (
	"time"
	"gorm.io/gorm"
)

type EmailVerification struct {
	gorm.Model
	Email       string    `json:"email" gorm:"column:email;unique"`
	CodeEmail   string    `json:"code_email" gorm:"column:codeEmail;"`
	ExpiredTime time.Time `json:"expired_time" gorm:"column:expiredTime;"`
	VerifyEmail bool      `json:"verify_email" gorm:"column:verifyEmail;"`
}

func (EmailVerification) TableName() string {
	return "email_verifications"
}
func (e *EmailVerification) SetDefaultValues() {
	e.VerifyEmail = false
}