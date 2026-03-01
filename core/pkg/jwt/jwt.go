package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateAccessToken 生成AccessToken
func GenerateAccessToken(uid int64, role string) (accessToken string, err error) {
	now := time.Now()
	// AccessToken
	accessClaims := MainClaims{
		Uid:  uid,
		Role: role,
		Type: JWT_TYPE_ACCESS_TOKEN,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   strconv.FormatInt(uid, 10),
			Audience:  []string{"user"},
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTTL)),
			NotBefore: jwt.NewNumericDate(now.Add(-5 * time.Second)), // 时钟偏差
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	// 生成未签名的Token
	accessTok := jwt.NewWithClaims(SigningMethod, accessClaims)
	accessToken, err = accessTok.SignedString(AccessSecret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %v", err)
	}
	// 返回
	return accessToken, nil
}

// GenerateRefreshToken 生成RefreshToken
func GenerateRefreshToken(uid int64, role string) (refreshToken string, err error) {
	now := time.Now()
	// RefreshToken
	refreshClaims := MainClaims{
		Uid:  uid,
		Role: role,
		Type: JWT_TYPE_REFRESH_TOKEN,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   strconv.FormatInt(uid, 10),
			Audience:  []string{"user"},
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTTL)),
			NotBefore: jwt.NewNumericDate(now.Add(-5 * time.Second)), // 时钟偏差
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	// 生成未签名的Token
	refreshTok := jwt.NewWithClaims(SigningMethod, refreshClaims)
	refreshToken, err = refreshTok.SignedString(RefreshSecret)
	if err != nil {
		return "", fmt.Errorf("sign refresh token: %v", err)
	}
	// 返回
	return refreshToken, nil
}

// VerifyAccessToken 校验AccessToken
func VerifyAccessToken(tokenStr string) (*MainClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is empty")
	}
	// 解析JWT token
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&MainClaims{},
		func(t *jwt.Token) (interface{}, error) {
			// 只接受 HS256 签名方法
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != SigningMethod.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return AccessSecret, nil
		},
		jwt.WithLeeway(5*time.Second),
	)
	// 错误判断
	if err != nil {
		// 使用 v5 的错误处理方式
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("malformed token")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token not valid yet")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("invalid signature")
		}

		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	// 类型断言
	claims, ok := token.Claims.(*MainClaims)
	if !ok {
		return nil, errors.New("invalid token claims structure")
	}
	// 检查token是否有效
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	// 检查token类型
	if claims.Type != JWT_TYPE_ACCESS_TOKEN {
		return nil, errors.New("token type mismatch: not an access token")
	}
	return claims, nil
}

// VerifyRefreshToken 校验AccessToken
func VerifyRefreshToken(tokenStr string) (*MainClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is empty")
	}
	// 解析JWT token
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&MainClaims{},
		func(t *jwt.Token) (interface{}, error) {
			// 只接受 HS256 签名方法
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != SigningMethod.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return RefreshSecret, nil
		},
		jwt.WithLeeway(5*time.Second),
	)
	// 错误判断
	if err != nil {
		// 使用 v5 的错误处理方式
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("malformed token")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token not valid yet")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("invalid signature")
		}

		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	// 类型断言
	claims, ok := token.Claims.(*MainClaims)
	if !ok {
		return nil, errors.New("invalid token claims structure")
	}
	// 检查token是否有效
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	// 检查token类型
	if claims.Type != JWT_TYPE_REFRESH_TOKEN {
		return nil, errors.New("token type mismatch: not an access token")
	}
	return claims, nil
}

func StripBearer(token string) string {
	if len(token) < 7 {
		return token
	}
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return token[7:] // 移除 "Bearer "
	}
	if strings.HasPrefix(strings.ToLower(token), "bearer") && len(token) > 6 {
		return token[6:] // 移除 "Bearer"
	}

	return token
}

func GetAccessTokenExpireTime() time.Duration {
	return AccessTTL
}

func GetRefreshTokenExpireTime() time.Duration {
	return RefreshTTL
}
