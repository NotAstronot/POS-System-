package rest

import (
	"context"
	"net/http"
	"pos-system/internal/domain"
	"pos-system/internal/usecase"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type loginAttempt struct {
	count    int
	firstTry time.Time
}

type AuthHandler struct {
	usecase   *usecase.AuthUsecase
	secretKey string
	mu        sync.Mutex
	attempts  map[string]*loginAttempt
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

	user, err := h.usecase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.mu.Lock()
		entry.count++
		h.mu.Unlock()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	h.mu.Lock()
	delete(h.attempts, ip)
	h.mu.Unlock()

	now := time.Now()

	// Tenant comes from the authenticated user (never from the client body).
	tenantID := int64(0)
	if user.TenantID != "" {
		if id, err := strconv.ParseInt(user.TenantID, 10, 64); err == nil {
			tenantID = id
		}
	}

	// Outlet: prefer outlet selected at login, fall back to the user's outlet.
	outletID := user.OutletID
	if req.OutletID != "" {
		outletID = req.OutletID
	}

	permissions, _ := h.usecase.GetPermissions(c.Request.Context(), tenantID, user.ID)
	permNames := make([]string, 0, len(permissions))
	for _, p := range permissions {
		permNames = append(permNames, p.Name)
	}

	// JWT payload: user_id, tenant_id, role, outlet_id (+ standard claims).
	claims := NewCustomClaims(user.ID, tenantID, user.Role, outletID, user.Name, now)
	tokenStr, err := SignToken(claims, h.secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token": tokenStr,
			"claims": gin.H{
				"user_id":   claims.UserID,
				"tenant_id": claims.TenantID,
				"role":      claims.Role,
				"outlet_id": claims.OutletID,
			},
			"user": gin.H{
				"id":          user.ID,
				"name":        user.Name,
				"role":        user.Role,
				"outletId":    outletID,
				"tenantId":    tenantID,
				"permissions": permNames,
			},
		},
	})
}
