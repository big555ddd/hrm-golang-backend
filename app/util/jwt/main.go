package jwt

import (
	"app/internal/logger"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const (
	VAL_USER  = "AUTH_USER"
	VAL_TOKEN = "AUTH_TOKEN"
	VAL_AGENT = "AUTH_AGENT"
)

type ClaimData struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Claims struct {
	Data ClaimData `json:"data"`
	Uuid string    `json:"uuid"`
	jwt.RegisteredClaims
}

func CreateToken(claims ClaimData) (string, error) {

	now := time.Now()
	id := uuid.New().String()
	duration := viper.GetInt64("JWT_DURATION")

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, Claims{
		claims,
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(duration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	})
	secret := []byte(viper.GetString("JWT_SECRET"))

	tokenString, err := token.SignedString(secret)
	if err != nil {
		logger.Errf("%s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func Verify(rawToken string) (*Claims, bool, error) {
	token, err := jwt.ParseWithClaims(rawToken, &Claims{}, getSecret)
	if err != nil {
		return nil, false, err
	}

	claims, ok := token.Claims.(*Claims)
	return claims, token.Valid && ok, nil
}

func GetClaims(c *gin.Context) (*ClaimData, error) {
	val, exists := c.Get(VAL_USER)
	if !exists {
		return nil, errors.New("claims doesn't exists")
	}
	data := val.(*ClaimData)
	return data, nil
}

func getSecret(token *jwt.Token) (interface{}, error) {

	return []byte(viper.GetString("JWT_SECRET")), nil
}

func GenerateExpires() time.Time {
	now := time.Now()
	duration := viper.GetInt64("JWT_DURATION")
	return now.Add(time.Duration(duration) * time.Hour)
}
