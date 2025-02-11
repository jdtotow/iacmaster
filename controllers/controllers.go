package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type DBController struct {
	dbUri      string
	dbPort     int
	dbUsername string
	dbPassword string
	dbName     string
}

func CreateDBController(uri, username, password, dbname string, port int) *DBController {
	return &DBController{
		dbUri:      uri,
		dbPort:     port,
		dbUsername: username,
		dbPassword: password,
		dbName:     dbname,
	}
}

func (db *DBController) Connect() error {
	err := errors.New("could not connect to DB")
	return err
}

func (db *DBController) Handle(context *gin.Context, concernedObject string) {
	fmt.Println("DBController reached, object called : ", concernedObject)
}

type SecurityController struct {
	secretKey []byte
}

func CreateSecurityController(secretKey string) *SecurityController {
	var secret []byte = []byte(secretKey)
	return &SecurityController{
		secretKey: secret,
	}
}

func (s *SecurityController) GetRole(username string) string {
	return ""
}

func (s *SecurityController) CreateToken(username string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "iacmaster",                      // Issuer
		"aud": s.GetRole(username),              // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}
	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}

func (s *SecurityController) VerifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	// Check for verification errors
	if err != nil {
		return nil, err
	}
	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	// Return the verified token
	return token, nil
}
