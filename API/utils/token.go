package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("PiCloudSecretKey")

type Claims struct {
    UserID int    `json:"user_id"`
    Username   string `json:"username"`
    Email  string `json:"email"` 
    jwt.StandardClaims
}


func GenerateJWT(userID int, username string, email string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        Username:   username,
        Email:  email, 
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

