package rest

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"pos-system/internal/database"
	"pos-system/internal/domain"
)

// TenantContext is the Tenant Extractor & RLS Injector middleware.
//
// Flow per request:
//  1. Read Bearer JWT from the Authorization header
//  2. Extract tenant_id (and user identity) from custom claims
//  3. Validate tenant_id (> 0, exists, status = active)
//  4. Inject identity into context.Context (domain) and gin context
//  5. Inject app.current_tenant into the PostgreSQL session (dedicated conn)
//
// Run after JWTAuth (or alone — it parses the token itself).
func TenantContext(secretKey string, db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- 1. Read JWT ---
		tokenString, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		// --- 2. Extract tenant_id from claims ---
		claims, err := ParseToken(tokenString, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		tenantID := claims.TenantID

		// --- 3. Validate tenant_id ---
		if tenantID <= 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid tenant"})
			return
		}
		if err := validateTenant(c.Request.Context(), db, tenantID); err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "tenant not found"})
			case errors.Is(err, errTenantInactive):
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "tenant is not active"})
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "tenant validation failed"})
			}
			return
		}

		// --- 4. Inject into context.Context (Go) ---
		ctx := c.Request.Context()
		ctx = domain.WithTenantID(ctx, tenantID)
		ctx = domain.WithUserID(ctx, claims.UserID)
		ctx = domain.WithUserRole(ctx, claims.Role)

		// --- 5. Inject into PostgreSQL session (RLS) ---
		conn, err := database.InjectRLSSession(ctx, db, tenantID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to start tenant database session"})
			return
		}
		// Released after all downstream handlers finish (c.Next returns).
		defer conn.Close()

		ctx = database.WithRLSConn(ctx, conn)
		c.Request = c.Request.WithContext(ctx)

		// Mirror onto gin keys for handlers that use getTenantID(c).
		c.Set("tenant_id", tenantID)
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("outlet_id", claims.OutletID)

		c.Next()
	}
}

// errTenantInactive indicates the tenant row exists but status != active.
var errTenantInactive = errors.New("tenant is not active")

// validateTenant checks that the tenant exists and is active.
// The tenants table has no RLS, so a direct query is safe.
func validateTenant(ctx context.Context, db *sql.DB, tenantID int64) error {
	var status string
	err := db.QueryRowContext(ctx, `SELECT status FROM tenants WHERE id = $1`, tenantID).Scan(&status)
	if err != nil {
		return err
	}
	if status != "active" {
		return errTenantInactive
	}
	return nil
}

// bearerToken extracts the raw JWT from an Authorization header value.
func bearerToken(header string) (string, bool) {
	if header == "" {
		return "", false
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", false
	}
	return parts[1], true
}
