package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
    UserID int    `json:"user_id"`
    Username   string `json:"username"`
    Email  string `json:"email"` 
    jwt.StandardClaims
}
