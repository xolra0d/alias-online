package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetEnvOrExit tries to get `name` env var. If not set - exits.
func GetEnvOrExit(name string) string {
	out := strings.TrimSpace(os.Getenv(name))
	if out == "" {
		fmt.Printf("ENV: `%s` is not specified!\n", name)
		os.Exit(1)
	}
	return out
}

// GetEnvOrFallback tries to get `name` env var and return it. If not set - returns `fallback`.
func GetEnvOrFallback(name, fallback string) string {
	out := strings.TrimSpace(os.Getenv(name))
	if out == "" {
		return fallback
	}
	return out
}

// StringToSeconds parses `val` (unsigned num) into num seconds. Uses `name` for panic.
func StringToSeconds(name, val string) time.Duration {
	out, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("ENV: `%s`=%s is invalid: %s!", name, val, err))
	}

	return time.Duration(out) * time.Second
}

// StringToUInt parses `val` (unsigned num) into uint64. Uses `name` for panic.
func StringToUInt(name, val string) uint64 {
	out, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("ENV: `%s`=%s is invalid: %s!", name, val, err))
	}

	return out
}

// StringToBool parses `val` into bool. Uses `name` for panic.
func StringToBool(name, val string) bool {
	out, err := strconv.ParseBool(val)
	if err != nil {
		panic(fmt.Sprintf("ENV: `%s`=%s is invalid: %s!", name, val, err))
	}

	return out
}
