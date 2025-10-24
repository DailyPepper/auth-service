package handlers

import (
	"auth-service/internal/models"
	"auth-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistrHandler struct {
	registrService service.Registr
}

func NewRegistrHandler(registrService service.Registr) *RegistrHandler {
	return &RegistrHandler{registrService: registrService}
}

func (h *RegistrHandler) NewRegistrService(c *gin.Context) {
	var req models.Registr

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	user, err := h.registrService.Registration(c.Request.Context(), &req)

	if err != nil {
		switch {
		case contains(err.Error(), "email already exists"):
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already registered",
			})
		case contains(err.Error(), "validation"):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Validation failed: " + err.Error(),
			})
		case contains(err.Error(), "password too weak"):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Password does not meet requirements",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Registration failed: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user.ToProfile(),
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && contains(s[1:], substr))
}
