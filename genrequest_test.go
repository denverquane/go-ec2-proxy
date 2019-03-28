package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

func TestCreateJWTRequest(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["DiscordID"] = ""
	claims["Type"] = "t2.micro"
	claims["Duration"] = 1
	claims["Quantity"] = 1
	claims["Data"] = 1

	tokenString, _ := token.SignedString([]byte("SECRET HERE"))
	fmt.Println("Created token: " + tokenString)
}
