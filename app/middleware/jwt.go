package middleware

import (
	"fmt"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/common/errcode"
	"gin/common/response"
	"gin/config"
	"gin/pkg/lang"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

type Jwt struct {
	base.BaseMiddleware
}

// Handle jwt中间件
func (s Jwt) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		errCode := errcode.Unauthorized()
		if token == "" || token == "null" {
			response.Error(c, &errCode)
			return
		}

		data, err := s.Decode(token)
		if err != nil {
			errCode = errCode.WithMsg(err.Error())
			response.Error(c, &errCode)
			return
		}

		c.Set(ctxkey.UserIdKey, data["id"])
		c.Next()
	}
}

// Encode 生成token
func (s Jwt) Encode(id int64, accessExp int64) (string, int64, error) {
	var (
		now = time.Now()
	)

	if accessExp == 0 {
		accessExp = now.Add(time.Duration(config.Jwt{}.Exp) * time.Second).Unix()
	} else {
		accessExp = now.Add(time.Duration(accessExp) * time.Second).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": accessExp,
	})

	jwtToken, err := token.SignedString([]byte(config.Jwt{}.Key))

	return jwtToken, accessExp, err
}

// Decode 解析token
func (s Jwt) Decode(jwtToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(lang.T(s.Ctx, "middleware.jwt.unsupportedSignatureMethod", nil)+": %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("your_secret_key")
		return []byte(config.Jwt{}.Key), nil
	})

	if err != nil {
		return nil, fmt.Errorf(lang.T(s.Ctx, "middleware.jwt.TokenParseErr", nil)+": %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf(lang.T(s.Ctx, "middleware.jwt.InvalidToken", nil))
}

// WithRefresh 刷新token
func (s Jwt) WithRefresh(id, accessExp, refreshExp int64) (accessToken, refreshToken string, tExp, rExp int64, err error) {
	var (
		now = time.Now()
	)

	if accessExp == 0 {
		accessExp = now.Add(time.Duration(config.Jwt{}.Exp) * time.Second).Unix()
	} else {
		accessExp = now.Add(time.Duration(accessExp) * time.Second).Unix()
	}

	if refreshExp == 0 {
		refreshExp = now.Add(time.Duration(config.Jwt{}.RefreshExp) * time.Second).Unix()
	} else {
		refreshExp = now.Add(time.Duration(refreshExp) * time.Second).Unix()
	}

	// Access Token
	accessClaims := jwt.MapClaims{
		"id":  id,
		"exp": accessExp,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString([]byte(config.Jwt{}.Key))
	if err != nil {
		return "", "", 0, 0, err
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"id":  id,
		"exp": refreshExp,
		"typ": "refresh", // 标记为 refresh token
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = rt.SignedString([]byte(config.Jwt{}.Key))
	if err != nil {
		return "", "", 0, 0, err
	}

	return accessToken, refreshToken, accessExp, refreshExp, nil
}
