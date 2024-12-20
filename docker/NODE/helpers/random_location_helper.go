package helpers

import (
    "math/rand"
	"time"
)

// Função para gerar uma localização dinâmica aleatória
func GetRandomLocation() string {
    locations := []string{
        "Lisboa",
        "Porto",
        "Madrid",
        "Paris",
        "Berlin",
        "Roma",
        "Nova Iorque",
        "Tóquio",
        "São Paulo",
        "Sydney",
    }

    rand.Seed(time.Now().UnixNano()) // Semente baseada no tempo atual
    return locations[rand.Intn(len(locations))]
}