package models 

type File struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Size int64  `json:"size"`
}