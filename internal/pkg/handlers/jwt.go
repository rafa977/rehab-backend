package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Username   string `json:"username"`
	ID         uint   `json:"id"`
	Authorized string `json:"authorized"`
	jwt.RegisteredClaims
}

var sampleSecretKey = []byte("SecretYouShouldHide")

func GenerateJWT(username string, id uint) (string, time.Time, error) {

	expTime := time.Now().Add(time.Minute * 60)

	claims := &Claims{
		Username:   username,
		ID:         id,
		Authorized: "authorized",
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(claims)

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		fmt.Println("here is the error")
		return "", expTime, err
	}

	return tokenString, expTime, nil

}

func ValidateToken(reqToken string) (string, uint, error) {

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	username := claims.Username
	id := claims.ID
	fmt.Println(username)

	fmt.Println(err)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", 0, errors.New("Invalid Signature")
		}
		return "", 0, errors.New("Bad Request")
	}
	if !tkn.Valid {
		return "", 0, errors.New("Token is invalid")
	}

	return username, id, nil
}
