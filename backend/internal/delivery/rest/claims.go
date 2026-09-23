package rest

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims is the JWT payload issued on login and validated on every request.
// It carries multi-tenant identity: tenant_id, user_id, role, and outlet_id.
type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	TenantID int64  `json:"tenant_id"`
	Role     string `json:"role"`
	OutletID string `json:"outlet_id"`
	Name     string `json:"name,omitempty"`
	jwt.RegisteredClaims
}

const tokenIssuer = "pos-system"
const tokenTTL = 24 * time.Hour

// NewCustomClaims builds signed-ready claims for an authenticated user.
func NewCustomClaims(userID, tenantID int64, role, outletID, name string, now time.Time) *CustomClaims {
	return &CustomClaims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		OutletID: outletID,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tokenIssuer,
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}
}

// SignToken signs claims with HS256 using the given secret.
func SignToken(claims *CustomClaims, secretKey string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
}

// ParseToken validates signature/expiry and returns typed custom claims.
func ParseToken(tokenString, secretKey string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	_, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithLeeway(30*time.Second),
	)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
