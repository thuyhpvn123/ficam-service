package request

type EmailAuthenticationRequest struct {
	Email       string    `json:"email" validate:"required,email"`
	CodeEmail   string    `json:"code_email" validate:"required"`
}