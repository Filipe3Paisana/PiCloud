package utils

import (
    "time"
    "net/http"
    "strings"
    "errors"
    "github.com/golang-jwt/jwt/v4"


    "api/models"
)


var jwtKey = []byte("PiCloudSecretKey")



// Função para gerar um token JWT
func GenerateJWT(userID int, username string, email string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // Tempo de expiração do token
    claims := &models.Claims{
        UserID: userID,
        Username: username,
        Email: email, 
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(), // Definindo o tempo de expiração
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Criando o token
    return token.SignedString(jwtKey) // Retornando o token assinado
}

// Função para verificar o token JWT e extrair as informações do utilizador
func VerifyJWT(tokenString string) (*models.Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil // Retorna a chave secreta para verificar a assinatura
    })

    if err != nil {
        return nil, errors.New("token inválido")
    }

    if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
        return claims, nil 
    }

    return nil, errors.New("token inválido ou expirado")
}


func ExtractUserIDFromJWT(r *http.Request) (int, error) {
    authHeader := r.Header.Get("Authorization") 
    if authHeader == "" {
        return 0, errors.New("token não fornecido")
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ") 
    userClaims, err := VerifyJWT(tokenString) 
    if err != nil {
        return 0, errors.New("token inválido ou expirado")
    }

    return userClaims.UserID, nil 
}
