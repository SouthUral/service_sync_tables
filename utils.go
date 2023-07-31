package main

import (
	"log"
	"os"
	"strconv"
)

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}

func getEnvInt(key string) int {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}
		return val
	}
	return 0
}

func CopyMap(data StateStorage) StateStorage {
	copyMap := make(StateStorage)
	for key, value := range data {
		copyMap[key] = value
	}
	return copyMap
}
