package request

type EmailVerificationRequest struct {
	Email       string    `json:"email" validate:"required,email"`
}
