package config

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// Config 全局配置
type Config struct {
	App    AppConfig
	MySQL  MySQLConfig
	Redis  RedisConfig
	JWT    JWTConfig
	AI     AIServiceConfig
	Upload UploadConfig
}

type AppConfig struct {
	Env  string
	Port string
}

type MySQLConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
}

type AIServiceConfig struct {
	URL   string
	Token string
}

type UploadConfig struct {
	MaxSizeMB         int
	AllowedExtensions []string
	StoragePath       string
}

var C *Config

func Init() {
	// 加载 .env 文件
	loadEnvFile("../../.env")
	loadEnvFile("../.env")
	loadEnvFile(".env")

	jwtSecret := getEnv("JWT_SECRET", "")
	if isInsecureSecret(jwtSecret) {
		jwtSecret = randomSecret()
		fmt.Println("[Config] WARNING: JWT_SECRET is missing or insecure; using an ephemeral startup secret")
	}

	C = &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "local"),
			Port: getEnv("APP_PORT", "8000"),
		},
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnv("MYSQL_PORT", "3306"),
			Database: getEnv("MYSQL_DATABASE", "novel2script"),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_PASSWORD", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       2,
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
		},
		AI: AIServiceConfig{
			URL:   getEnv("AI_SERVICE_URL", "http://localhost:9000"),
			Token: getEnv("AI_SERVICE_TOKEN", ""),
		},
		Upload: UploadConfig{
			MaxSizeMB:         20,
			AllowedExtensions: []string{".txt", ".md", ".docx"},
			StoragePath:       getEnv("STORAGE_PATH", "./storage"),
		},
	}

	fmt.Printf("[Config] MySQL: %s:%s@%s:%s/%s\n", C.MySQL.User, "***", C.MySQL.Host, C.MySQL.Port, C.MySQL.Database)
}

func isInsecureSecret(secret string) bool {
	trimmed := strings.TrimSpace(secret)
	if trimmed == "" {
		return true
	}
	insecureValues := map[string]bool{
		"default_secret_change_me":                            true,
		"change_me":                                           true,
		"change_me_to_a_long_random_string":                   true,
		"change_me_to_a_long_random_string_at_least_32_chars": true,
	}
	return insecureValues[trimmed] || len(trimmed) < 32
}

func randomSecret() string {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("ephemeral-%d", os.Getpid())
	}
	return hex.EncodeToString(buf)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// loadEnvFile 从文件加载环境变量
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 去除引号
		value = strings.Trim(value, "\"'")

		// 只在环境变量未设置时才设置
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}
