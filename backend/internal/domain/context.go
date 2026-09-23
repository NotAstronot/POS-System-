package domain

import "context"

type contextKey string

const TenantIDKey contextKey = "tenant_id"
const UserIDKey contextKey = "user_id"
const UserRoleKey contextKey = "user_role"

func WithTenantID(ctx context.Context, tenantID int64) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

func GetTenantID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(TenantIDKey).(int64)
	return id, ok
}

func MustGetTenantID(ctx context.Context) int64 {
	id, ok := GetTenantID(ctx)
	if !ok {
		panic("tenant_id not found in context")
	}
	return id
}

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func GetUserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserIDKey).(int64)
	return id, ok
}

func WithUserRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, UserRoleKey, role)
}

func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}
