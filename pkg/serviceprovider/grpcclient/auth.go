package grpcclient

import (
	"context"
	"errors"
	"gin/common/ctxkey"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// MethodAuthService 方法鉴权接口
type MethodAuthService interface {
	AuthMethods() map[string]bool // RPC方法名 -> 是否需要鉴权
}

var (
	authMethodServices = make(map[string]map[string]bool)
	authServicesMu     sync.RWMutex
	jwtKey             string
)

// RegisterAuth 注册grpc服务方法鉴权配置
func RegisterAuth(serviceName string, service MethodAuthService) {
	authServicesMu.Lock()
	defer authServicesMu.Unlock()
	authMethodServices[serviceName] = service.AuthMethods()
}

// SetJwtKey 设置JWT密钥
func SetJwtKey(key string) {
	jwtKey = key
}

// WithToken 在上下文中携带Token
func WithToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

// authUnaryServerInterceptor 服务端鉴权拦截器
func authUnaryServerInterceptor(ctx context.Context, req any, info *grpclib.UnaryServerInfo, handler grpclib.UnaryHandler) (any, error) {
	ctx, err := authContext(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

// authContext 校验方法鉴权并返回带用户ID的上下文
func authContext(ctx context.Context, fullMethod string) (context.Context, error) {
	serviceName, methodName := splitFullMethod(fullMethod)
	if !methodAuthRequired(serviceName, methodName) {
		return ctx, nil
	}

	token := authToken(ctx)
	if token == "" {
		return ctx, status.Error(codes.Unauthenticated, "Token不存在")
	}

	claims, err := decodeAuthToken(token)
	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, "Token无效或已过期")
	}

	id, ok := claims["id"].(float64)
	if !ok || int64(id) <= 0 {
		return ctx, status.Error(codes.Unauthenticated, "Token无效")
	}
	return ctxkey.WithValue(ctx, ctxkey.UserIdKey, int64(id)), nil
}

// splitFullMethod 从完整方法名提取服务名和方法名
func splitFullMethod(fullMethod string) (string, string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/")
	if idx := strings.LastIndex(fullMethod, "/"); idx > 0 {
		return fullMethod[:idx], fullMethod[idx+1:]
	}
	return fullMethod, ""
}

// methodAuthRequired 判断方法是否需要鉴权,未配置的方法默认不鉴权
func methodAuthRequired(serviceName, methodName string) bool {
	authServicesMu.RLock()
	defer authServicesMu.RUnlock()
	if methods, ok := authMethodServices[serviceName]; ok {
		return methods[methodName]
	}
	return false
}

// authToken 从metadata获取Token
func authToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if values := md.Get("authorization"); len(values) > 0 {
		return strings.TrimSpace(strings.TrimPrefix(values[0], "Bearer "))
	}
	if values := md.Get("token"); len(values) > 0 {
		return strings.TrimSpace(values[0])
	}
	return ""
}

// decodeAuthToken 解析JWT
func decodeAuthToken(token string) (map[string]any, error) {
	if jwtKey == "" {
		return nil, errors.New("jwt key not configured")
	}

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unsupported signing method")
		}
		return []byte(jwtKey), nil
	})
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
