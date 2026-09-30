package ws

import "context"

type userIDKey struct{}

// WithUserID 注入登录用户ID
func WithUserID(ctx context.Context, userID int64) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, userIDKey{}, userID)
}

// UserID 获取登录用户ID
func UserID(ctx context.Context) int64 {
	if ctx == nil {
		return 0
	}
	userID, ok := ctx.Value(userIDKey{}).(int64)
	if !ok || userID <= 0 {
		return 0
	}
	return userID
}
