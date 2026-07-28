package rest

import (
	"net/http"
	"pos-system/internal/usecase"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type loginAttempt struct {
	count    int
	firstTry time.Time
}

type AuthHandler struct {
	usecase    *usecase.AuthUsecase
	secretKey  string
	mu         sync.Mutex
	attempts   map[string]*loginAttempt
}

func NewAuthHandler(usecase *usecase.AuthUsecase, secretKey string) *AuthHandler {
	return &AuthHandler{
		usecase:   usecase,
		secretKey: secretKey,
		attempts:  make(map[string]*loginAttempt),
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	TenantID string `json:"tenant_id"`
	OutletID string `json:"outlet_id"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Rate limiting: max 5 attempts per IP per 15 minutes
	ip := c.ClientIP()
	h.mu.Lock()
	entry, exists := h.attempts[ip]
	if !exists {
		entry = &loginAttempt{firstTry: time.Now()}
		h.attempts[ip] = entry
	}
	h.mu.Unlock()
	if time.Since(entry.firstTry) > 15*time.Minute {
		h.mu.Lock()
		delete(h.attempts, ip)
		h.mu.Unlock()
	} else if entry.count >= 5 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many login attempts, try again later"})
		return
	}

	user, err := h.usecase.Login(req.Username, req.Password)
	if err != nil {
		h.mu.Lock()
		entry.count++
		h.mu.Unlock()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Success — reset rate limit
	h.mu.Lock()
	delete(h.attempts, ip)
	h.mu.Unlock()

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"name":  user.Name,
		"role":  user.Role,
		"iat":   now.Unix(),
		"exp":   now.Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token": tokenStr,
			"user": gin.H{
				"id":       user.ID,
				"name":     user.Name,
				"role":     user.Role,
				"outletId": user.OutletID,
				"tenantId": user.TenantID,
			},
		},
	})
}
