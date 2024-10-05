package models

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password_hash"`
    Email    string `json:"email"`
}
