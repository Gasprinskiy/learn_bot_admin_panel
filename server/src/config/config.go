package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config структура для хранения переменных окружения
type Config struct {
	PostgresURL  string
	TgApiURL     string
	BotToken     string
	RedisAddr    string
	RedisPass    string
	JwtSecret    string
	ServerPort   string
	AllowOrigins []string
	SSETTL       time.Duration
	JwtSecretTTL time.Duration
	RedisTTL     time.Duration
	GRPCAddr     string
}

// NewConfig загружает переменные из .env и возвращает структуру Config
func NewConfig() *Config {
	redisTtl, err := strconv.Atoi(os.Getenv("REDIS_TTL"))
	if err != nil {
		log.Panic("не удалось получить время жизни кеша: ", err)
	}

	jwtSecretTtl, err := strconv.Atoi(os.Getenv("JWT_SECRET_TTL"))
	if err != nil {
		log.Panic("не удалось получить время жизни jwt токена: ", err)
	}

	sseTtl, err := strconv.Atoi(os.Getenv("SSE_TTL"))
	if err != nil {
		log.Panic("не удалось получить время жизни sse канала: ", err)
	}

	var allowOrigins []string

	allowOriginsStr := os.Getenv("ALLOW_ORIGINS")
	if allowOriginsStr == "" {
		allowOrigins = append(allowOrigins, "*")
	} else {
		for _, origin := range strings.Split(allowOriginsStr, ",") {
			trimed := strings.TrimSpace(origin)
			allowOrigins = append(allowOrigins, trimed)
		}
	}

	return &Config{
		PostgresURL:  os.Getenv("POSTGRES_URL"),
		TgApiURL:     os.Getenv("TG_API_URL"),
		BotToken:     os.Getenv("BOT_TOKEN"),
		RedisPass:    os.Getenv("REDIS_PASSWORD"),
		JwtSecret:    os.Getenv("JWT_SECRET"),
		GRPCAddr:     os.Getenv("GRPC_ADDR"),
		AllowOrigins: allowOrigins,
		ServerPort:   fmt.Sprintf(":%s", os.Getenv("HTTP_SERVER_PORT")),
		RedisAddr:    fmt.Sprintf("redis:%s", os.Getenv("REDIS_PORT")),
		RedisTTL:     time.Minute * time.Duration(redisTtl),
		JwtSecretTTL: time.Hour * time.Duration(jwtSecretTtl),
		SSETTL:       time.Minute * time.Duration(sseTtl),
	}
}
