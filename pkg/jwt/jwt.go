package jwt

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lyric1024/blog-system/pkg/errors"
)

var (
	jwtService *JWTService
	once       sync.Once
)

type JWTService struct {
	Secret     []byte
	ExpireTime uint
}

type Claims struct {
	UserID uint `json:"userID"`
	jwt.RegisteredClaims
}

func Init(secret string, expireTime uint) {
	if secret == "" {
		panic("JWT secret cannot be empty")
	}
	if expireTime <= 0 {
		expireTime = 168
	}

	once.Do(func() {
		jwtService = &JWTService{
			Secret:     []byte(secret),
			ExpireTime: expireTime,
		}
	})
}

// 获取用户token
func GetToken(userID uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(jwtService.ExpireTime) * time.Hour)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtService.Secret)
}

// 解析token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtService.Secret, nil
	})

	if err != nil {
		return nil, errors.Internal("解析token失败", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.BadRequest("无效的token", nil)
}

// 将token写入cookies
func SetTokenCookie(c *gin.Context, token string) {
	maxAge := int(jwtService.ExpireTime * 3600)
	c.SetCookie(
		"token",
		token,
		maxAge,
		"/",
		"",
		false,
		true,
	)
}
