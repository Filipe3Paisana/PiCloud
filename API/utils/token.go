package utils

import (
    "time"
    "net/http"
    "strings"
    "errors"
    "github.com/golang-jwt/jwt/v4"
)

// Chave secreta utilizada para assinar os tokens JWT
var jwtKey = []byte("PiCloudSecretKey")

// Estrutura que representa os claims do JWT
type Claims struct {
    UserID   int    `json:"user_id"` // ID do usuário
    Username string `json:"username"`
    Email    string `json:"email"` 
    jwt.StandardClaims
}

// Função para gerar um token JWT
func GenerateJWT(userID int, username string, email string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // Tempo de expiração do token
    claims := &Claims{
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

// Função para verificar o token JWT e extrair as informações do usuário
func VerifyJWT(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil // Retorna a chave secreta para verificar a assinatura
    })

    if err != nil {
        return nil, errors.New("token inválido")
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil // Retorna os claims se o token for válido
    }

    return nil, errors.New("token inválido ou expirado")
}

// Função para extrair o userID da requisição
func ExtractUserIDFromJWT(r *http.Request) (int, error) {
    authHeader := r.Header.Get("Authorization") // Obtém o cabeçalho Authorization
    if authHeader == "" {
        return 0, errors.New("token não fornecido")
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ") // Remove o prefixo "Bearer "
    userClaims, err := VerifyJWT(tokenString) // Verifica e extrai os claims do token
    if err != nil {
        return 0, errors.New("token inválido ou expirado")
    }

    return userClaims.UserID, nil // Retorna o userID extraído do token
}
