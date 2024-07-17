package env_config

import (
	"os"
	"strconv"
	"time"
)

// Env String
func Env(key, defaultValue string) (value string) {
	if value = os.Getenv(key); value == "" {
		value = defaultValue
	}
	return
}

// EnvBytes Byte
func EnvBytes(key, defaultValue string) []byte {
	return []byte(Env(key, defaultValue))
}

// EnvDuration time Duration
func EnvDuration(key, defaultValue string) (time.Duration, error) {
	return time.ParseDuration(Env(key, defaultValue))
}

// EnvInt Integer
func EnvInt(key, defaultValue string) (int, error) {
	return strconv.Atoi(Env(key, defaultValue))
}

func EnvInt8(key, defaultValue string) (int8, error) {
	value, err := strconv.ParseInt(Env(key, defaultValue), 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(value), nil
}

func EnvInt16(key, defaultValue string) (int16, error) {
	value, err := strconv.ParseInt(Env(key, defaultValue), 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(value), nil
}

func EnvInt32(key, defaultValue string) (int32, error) {
	value, err := strconv.ParseInt(Env(key, defaultValue), 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(value), nil
}

func EnvInt64(key, defaultValue string) (int64, error) {
	return strconv.ParseInt(Env(key, defaultValue), 10, 64)
}

// EnvUint Uint
func EnvUint(key, defaultValue string) (uint, error) {
	value, err := strconv.ParseUint(Env(key, defaultValue), 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(value), nil
}

func EnvUint8(key, defaultValue string) (uint8, error) {
	value, err := strconv.ParseUint(Env(key, defaultValue), 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(value), nil
}

func EnvUint16(key, defaultValue string) (uint16, error) {
	value, err := strconv.ParseUint(Env(key, defaultValue), 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(value), nil
}

func EnvUint32(key, defaultValue string) (uint32, error) {
	value, err := strconv.ParseUint(Env(key, defaultValue), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}

func EnvUint64(key, defaultValue string) (uint64, error) {
	return strconv.ParseUint(Env(key, defaultValue), 10, 64)
}

// EnvFloat32 Float
func EnvFloat32(key, defaultValue string) (float32, error) {
	value, err := strconv.ParseFloat(Env(key, defaultValue), 32)
	if err != nil {
		return 0, err
	}
	return float32(value), nil
}

func EnvFloat64(key, defaultValue string) (float64, error) {
	return strconv.ParseFloat(Env(key, defaultValue), 64)
}

// EnvBool Boolean
func EnvBool(key, defaultValue string) (bool, error) {
	return strconv.ParseBool(Env(key, defaultValue))
}
