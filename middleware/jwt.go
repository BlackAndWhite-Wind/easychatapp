package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var TokenExpired = errors.New("token is expired")
var jwtSecret = []byte("ice_moss")

type Claims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

func JWY() gin.HandlerFunc {
	return func(c *gin.Context) {
		//tokens := c.PostForm("token")
		token := c.Query("token")
		user := c.Query("userId")
		userId, err := strconv.Atoi(user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "您userId不合法",
			})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "请登录",
			})
			c.Abort()
			return
		} else {
			claims, err := ParseToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "token失效",
				})
				c.Abort()
				return
			} else if time.Now().Unix() > claims.ExpiresAt {
				err = TokenExpired
				c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "授权已过期",
				})
				c.Abort()
				return
			}

			if claims.UserID != uint(userId) {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "您的登录不合法",
				})
				c.Abort()
				return
			}
			fmt.Println("token认证成功")
			c.Next()
		}
	}
}

// GenerateToken 根据用户的用户名和密码产生token
func GenerateToken(userId uint, iss string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(48 * 30 * time.Hour)
	claims := Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    iss,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
