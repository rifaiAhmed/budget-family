package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Migrate  MigrateConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Mode    string
	Address string
}

type DatabaseConfig struct {
	DSN string
}

type RedisConfig struct {
	Enabled bool
	Addr    string
	Pass    string
	DB      int
}

type AuthConfig struct {
	JWTSecret        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	PasswordCost     int
	Issuer           string
	RefreshTokenName string
}

type MigrateConfig struct {
	Auto bool
	Path string
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() Config {
	// Optional .env support for local development.
	// If .env doesn't exist, we silently continue.
	viper.SetConfigFile(".env")
	if _, err := os.Stat(".env"); err == nil {
		_ = viper.ReadInConfig()
	}

	viper.SetDefault("SERVER_MODE", "release")
	viper.SetDefault("SERVER_ADDRESS", ":8080")

	viper.SetDefault("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=budget_family port=5432 sslmode=disable")

	viper.SetDefault("REDIS_ENABLED", false)
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASS", "")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("AUTH_JWT_SECRET", "change-me")
	viper.SetDefault("AUTH_ACCESS_TTL", "15m")
	viper.SetDefault("AUTH_REFRESH_TTL", "720h")
	viper.SetDefault("AUTH_PASSWORD_COST", 12)
	viper.SetDefault("AUTH_ISSUER", "budget-family")

	viper.SetDefault("MIGRATE_AUTO", true)
	viper.SetDefault("MIGRATE_PATH", "file://migrations")

	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:5173")

	viper.AutomaticEnv()

	accessTTL, _ := time.ParseDuration(viper.GetString("AUTH_ACCESS_TTL"))
	refreshTTL, _ := time.ParseDuration(viper.GetString("AUTH_REFRESH_TTL"))

	return Config{
		Server: ServerConfig{
			Mode:    viper.GetString("SERVER_MODE"),
			Address: viper.GetString("SERVER_ADDRESS"),
		},
		Database: DatabaseConfig{
			DSN: viper.GetString("DATABASE_DSN"),
		},
		Redis: RedisConfig{
			Enabled: viper.GetBool("REDIS_ENABLED"),
			Addr:    viper.GetString("REDIS_ADDR"),
			Pass:    viper.GetString("REDIS_PASS"),
			DB:      viper.GetInt("REDIS_DB"),
		},
		Auth: AuthConfig{
			JWTSecret:       viper.GetString("AUTH_JWT_SECRET"),
			AccessTokenTTL:  accessTTL,
			RefreshTokenTTL: refreshTTL,
			PasswordCost:    viper.GetInt("AUTH_PASSWORD_COST"),
			Issuer:          viper.GetString("AUTH_ISSUER"),
		},
		Migrate: MigrateConfig{
			Auto: viper.GetBool("MIGRATE_AUTO"),
			Path: viper.GetString("MIGRATE_PATH"),
		},
		CORS: CORSConfig{
			AllowedOrigins: splitCSV(viper.GetString("CORS_ALLOWED_ORIGINS")),
		},
	}
}

func splitCSV(s string) []string {
	parts := []string{}
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		parts = append(parts, p)
	}
	return parts
}
