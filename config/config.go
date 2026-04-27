// config/config.go
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ServerPort string
}

func Load() *Config {
	envMap := readEnvFile(".env")
	return &Config{
		DBHost:     getEnv(envMap, "DB_HOST", "localhost"),
		DBPort:     getEnv(envMap, "DB_PORT", "3306"),
		DBUser:     getEnv(envMap, "DB_USER", "root"),
		DBPassword: getEnv(envMap, "DB_PASSWORD", "password"),
		DBName:     getEnv(envMap, "DB_NAME", "auth_db"),
		JWTSecret:  getEnv(envMap, "JWT_SECRET", "your-super-secret-key-change-in-production"),
		ServerPort: getEnv(envMap, "SERVER_PORT", "8080"),
	}
}

func readEnvFile(filename string) map[string]string {
	envMap := make(map[string]string)
	file, oErr := os.Open(filename)
	if oErr != nil {
		return envMap
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			envMap[key] = value
		}
	}
	return envMap
}
func getEnv(envMap map[string]string, key, dValue string) string {
	if value, exists := envMap[key]; exists {
		return value
	}
	return dValue
}
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
