package config

import "os"

type Config struct {
	DBName     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
}

// 取得环境变量
func LoadConfig() Config {
	return Config{
		DBName:     getEnv("MYSQL_DATABASE", "user_memo"),
		DBHost:     getEnv("DB_HOST", "mysql"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("MYSQL_USER", "root"),
		DBPassword: getEnv("MYSQL_PASSWORD", "Johntang2005"),
	}
}

// 解析环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
