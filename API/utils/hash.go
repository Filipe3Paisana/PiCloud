package utils

import (
    "golang.org/x/crypto/bcrypt"
)
// HashPassword recebe uma senha em texto puro e retorna um hash seguro
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

// ComparePassword compara uma senha em texto puro com um hash
func ComparePassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

