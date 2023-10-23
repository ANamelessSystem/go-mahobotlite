package main

import (
	"os"
	"strings"
)

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getBoolEnvWithDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.ToLower(value) == "true"
}

var (
	// Get param from env. If empty, use default.
	CqHttpAddr = getEnvWithDefault("CQ_ADDRESS", "http://127.0.0.1:5700")
	LsnrAddr   = getEnvWithDefault("LISTENER_ADDRESS", "0.0.0.0:8876")
	ShowLogs   = getBoolEnvWithDefault("SHOW_LOGS", true)
)
