package rest

import (
	"database/sql"
	"net/http"
	"pos-system/internal/database"
	"pos-system/internal/domain"

	"github.com/gin-gonic/gin"
)

func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := ParseToken(tokenString, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Expose custom claims to handlers (gin context + request context).
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("user_role", claims.Role)
		c.Set("outlet_id", claims.OutletID)

		ctx := c.Request.Context()
		ctx = domain.WithUserID(ctx, claims.UserID)
		ctx = domain.WithTenantID(ctx, claims.TenantID)
		ctx = domain.WithUserRole(ctx, claims.Role)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// Chain runs middlewares in order; stops early if one aborts the request.
func Chain(mws ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, mw := range mws {
			mw(c)
			if c.IsAborted() {
				return
			}
		}
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || (role != "admin" && role != "super_admin") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "akses ditolak: hanya admin"})
			return
		}
		c.Next()
	}
}

func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "super_admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "akses ditolak: hanya super admin"})
			return
		}
		c.Next()
	}
}

func RequirePermission(db *sql.DB, permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role == "admin" || role == "super_admin" {
			c.Next()
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "akses ditolak"})
			return
		}

		var count int
		count, err := database.WithTenantTx1(c.Request.Context(), db, func(tx *sql.Tx, tenantID int64) (int, error) {
			var n int
			qerr := tx.QueryRowContext(c.Request.Context(),
				`SELECT COUNT(*) FROM user_permissions up
				 INNER JOIN permissions p ON p.id = up.permission_id
				 INNER JOIN users u ON u.id = up.user_id
				 WHERE up.user_id = $1 AND p.name = $2 AND u.tenant_id = $3::text`,
				userID, permissionName, tenantID,
			).Scan(&n)
			return n, qerr
		})
		if err != nil || count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "akses ditolak: tidak ada hak akses"})
			return
		}
		c.Next()
	}
}
