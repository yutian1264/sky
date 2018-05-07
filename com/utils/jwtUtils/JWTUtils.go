package jwtUtils

import (
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"fmt"
	"log"
	"github.com/dgrijalva/jwt-go/request"
	"time"
)

type Token struct {
	Token string `json:"token"`
}

/**
	Secret key
 */
var SecretKey = "abcdefghijklmnopqrstuvwxyz"


func CreateToken()(error,string){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		log.Println("CreateToken", "Error while signing the token")
		fatal(err)
	}

	return err,tokenString
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}