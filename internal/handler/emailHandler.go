package handler

import (
	"meta-node-ficam/internal/service"
	"meta-node-ficam/utils"
	"meta-node-ficam/internal/request"
	"net/http"
	"github.com/gin-gonic/gin"
)

type EmailHandler interface {
	EmailVerification(c *gin.Context)
	EmailAuthentication(c *gin.Context)
}

type emailHandler struct {
	emailService service.EmailService
}

func NewEmailHandler(emailService service.EmailService) EmailHandler {
	return &emailHandler{emailService}
}

func (h *emailHandler) EmailVerification(c *gin.Context) {
	var request request.EmailVerificationRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := utils.ValidateStruct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.emailService.EmailVerification(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email Verification Initiated. Code sent to your email. "})
}

func (h *emailHandler) EmailAuthentication(c *gin.Context) {
	var request request.EmailAuthenticationRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request",
			"status": http.StatusBadRequest})
		return
	}
	if err := utils.ValidateStruct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),
			"status": http.StatusBadRequest})
		return
	}
	if err := h.emailService.EmailAuthentication(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(),
			"status": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email Authentication Success",
		"status": http.StatusOK})
}
