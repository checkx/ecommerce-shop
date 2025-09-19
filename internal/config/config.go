package config

import (
	"os"
	"strconv"
)

type Config struct {
	Env                   string
	HTTPAddr              string
	DBURL                 string
	JWTSecret             string
	ReservationTTLMinutes int
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() Config {
	minutesStr := getEnv("RESERVATION_TTL_MINUTES", "15")
	minutes, err := strconv.Atoi(minutesStr)
	if err != nil {
		minutes = 15
	}
	return Config{
		Env:                   getEnv("APP_ENV", "development"),
		HTTPAddr:              getEnv("HTTP_ADDR", ":8080"),
		DBURL:                 getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable"),
		JWTSecret:             getEnv("JWT_SECRET", "secret"),
		ReservationTTLMinutes: minutes,
	}
}
