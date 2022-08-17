package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type TokenEngine interface {
	Encode(user string, uniqueID string, duration time.Duration) (string, error)
	Decode(encoded string) (string, error)
}

type TokenEngineImpl struct {
	jwtSecret []byte
}

func NewTokenEngine(configuration Configuration) (TokenEngine, error) {

	// Get the JWT secret string, and convert it to the hex byte sequence
	if rawString, isPresent := configuration.GetConfigurationString("JUTZO_JWT_SECRET"); isPresent {
		engine := new(TokenEngineImpl)
		jwtSecret, err := hex.DecodeString(rawString)
		engine.jwtSecret = jwtSecret
		return engine, err
	} else {
		return nil, errors.New("could not get the JWT secret key from configuration")
	}
}

func (engine *TokenEngineImpl) Encode(user string, uniqueID string, duration time.Duration) (string, error) {

	// Create a JWT token that will give us the information we need to get back
	// to that user session. Effectively we just store the Id, which will point to the
	// Redis UserSession object, but we add the Subject, Issuer, and ExpiresAt as well
	// (our sessions are 8 hours)
	cookie := jwt.StandardClaims{
		Issuer:    "Jutzo Service",
		Subject:   user,
		Id:        uniqueID,
		ExpiresAt: time.Now().Local().Add(duration).Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cookie)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(engine.jwtSecret)

}

func (engine *TokenEngineImpl) Decode(encoded string) (string, error) {
	// Attempt to parse the token. This will give us a set of standard
	// claims, of which the Id will point us to the user session in Redis
	token, err := jwt.ParseWithClaims(encoded, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {

		// Validate that we got the right signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return engine.jwtSecret, nil
	})

	if err != nil {
		return "", err
	} else {

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if ok && token.Valid {
			return claims.Id, nil
		} else {
			return "", errors.New("unable to extract JWT claims")
		}
	}

}
