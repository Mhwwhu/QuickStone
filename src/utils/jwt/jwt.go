package jwt

import (
	"QuickStone/src/common"
	"QuickStone/src/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func init() {
	jwtKey = []byte(config.EnvCfg.JwtSecretKey)
}

type Claims struct {
	UserID   common.UserIdT `json:"user_id"`
	Username string         `json:"username"`
	jwt.RegisteredClaims
}

func GetToken(userId common.UserIdT, userName string) string {
	claims := &Claims{
		UserID:   userId,
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                    // 签发时间
			Issuer:    "QuickStone",                                      // 签发者
		},
	}

	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func VerifyToken(token string) (Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return Claims{}, err
	}

	if claims, ok := parsedToken.Claims.(Claims); ok && parsedToken.Valid {
		return claims, nil
	}
	return Claims{}, err
}
