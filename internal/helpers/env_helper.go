package helpers

import "os"

// GetEnvOrDefault
//получает: имя переменной окружения и значение по умолчанию
//возвращает: значение из переменной окружения или значение по умолчанию
func GetEnvOrDefault(key, defaultValue string) string {
	var value = defaultValue
	if valueFromEnv := os.Getenv(key); valueFromEnv != "" {
		value = valueFromEnv
	}
	return value
}
