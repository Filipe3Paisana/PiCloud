package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)


var jwtKey = []byte("PiCloudSecretKey")


func GenerateJWT(userID int) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // 
    claims := &jwt.StandardClaims{
        Subject:   string(userID), 
        ExpiresAt: expirationTime.Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}


func ValidateJWT(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
}
