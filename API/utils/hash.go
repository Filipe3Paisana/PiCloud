package utils

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword recebe uma senha em texto puro e retorna um hash seguro
func HashPassword(password string) (string, error) {
    // Gera o hash da senha com o custo padrão do bcrypt
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func ComparePassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
