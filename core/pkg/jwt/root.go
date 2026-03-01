package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	SigningMethod = jwt.SigningMethodHS256 // 签名方法
	AccessSecret  = []byte("leave_blog_jwt_access_secret")
	RefreshSecret = []byte("leave_blog_video_jwt_refresh_secret")
	Issuer        = "leave.blog"
	AccessTTL     = 24 * time.Hour     // 访问令牌有效期
	RefreshTTL    = 7 * 24 * time.Hour // 刷新令牌有效期

	JWT_TYPE_ACCESS_TOKEN  = "access"
	JWT_TYPE_REFRESH_TOKEN = "refresh"
	JWT_CONTEXT_KEY        = "jwt_context_key"
	JWT_REFRESH_KEY        = "jwt_refresh_key"
	JWT_ROLE_ADMIN         = "jwt_role_admin"
	JWT_ROLE_USER          = "jwt_role_user"
)

type MainClaims struct {
	Uid  int64  `json:"uid"`
	Role string `json:"role"`
	Type string `json:"type"`

	jwt.RegisteredClaims
}
