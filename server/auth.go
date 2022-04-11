package server

import (
	"crypto/sha1"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/karthiklsarma/cedar-logging/logging"
	"github.com/karthiklsarma/cedar-schema/gen"
	"github.com/karthiklsarma/cedar-storage/storage"
)

const (
	TOKEN_KEY          = "TOKEN_KEY"
	USER_TABLE_CONNSTR = "USER_TABLE_CONNSTR"
)

func getSigningKey() []byte {
	return []byte(os.Getenv(TOKEN_KEY))
}

func getUserTableConnStr() string {
	return os.Getenv(USER_TABLE_CONNSTR)
}

type CustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func RegisterUser(user *gen.User) (bool, error) {
	var storageSink storage.IStorageSink = &storage.CosmosSink{}

	if err := storageSink.Connect(); err != nil {
		logging.Error(fmt.Sprintf("Failed to connect to storage sink: %v", err))
		return false, err
	}

	sha := sha1.New()
	sha.Write([]byte(user.Password))
	user.Password = fmt.Sprintf("%x", sha.Sum(nil))
	if status, err := storageSink.InsertUser(user); !status || err != nil {
		logging.Error(fmt.Sprintf("Unable to register user. Error : %v", err))
		return status, err
	}

	return true, nil
}

func AuthenticateUser(username, password string) (bool, error) {
	var storageSink storage.IStorageSink = &storage.CosmosSink{}
	if err := storageSink.Connect(); err != nil {
		logging.Error(fmt.Sprintf("Failed to connect to storage sink: %v", err))
		return false, err
	}

	sha := sha1.New()
	sha.Write([]byte(password))
	hashPass := sha.Sum(nil)
	status, err := storageSink.Authenticate(username, fmt.Sprintf("%x", hashPass))
	if err != nil {
		logging.Error(fmt.Sprintf("Unable to authenticate user. Error : %v", err))
		return false, err
	}

	return status, nil
}

func GetNewToken(username string) string {
	claims := CustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer:    "cedar-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := getSigningKey()
	userToken, err := token.SignedString(key)
	if err != nil {
		logging.Error(fmt.Sprintf("error generating token for user: %v, error: %v", username, err))
	}

	return userToken
}

func ValidateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return getSigningKey(), nil
	})

	if err != nil {
		logging.Error(fmt.Sprintf("Error parsing token: %v", err))
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		logging.Info(fmt.Sprintf("User %v authenticated.", claims["username"]))
		return true
	}

	return false
}
