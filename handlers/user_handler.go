package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"vigilant-spork/models"
	"vigilant-spork/services"
)

type UserHandler struct {
	Service *services.UserService
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var signUp models.User
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&signUp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call Service layer
	err := h.Service.RegisterUser(&signUp)
	//Handle errors
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict for duplicate emails
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest) // other validation errors
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(signUp)
}

func (h *UserHandler) Logout(c *gin.Context) {
	// Get the raw token string from the context
	tokenValue, exists := c.Get("jwtTokenString")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
		return
	}

	tokenString := tokenValue.(string)

	// Call the service to revoke session (blacklist)
	err := h.Service.RevokeSession(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to revoke session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
