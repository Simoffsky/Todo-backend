package config

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvConfig(t *testing.T) {
	os.Setenv("LOGGER_LEVEL", "INFO")
	os.Setenv("HTTP_TIMEOUT", "10")
	os.Setenv("SERVER_PORT", "9090")

	config := NewEnvConfig()

	assert.Equal(t, "INFO", config.LoggerLevel)
	assert.Equal(t, 10*time.Second, config.HTTPTimeout)
	assert.Equal(t, "9090", config.ServerPort)

	os.Unsetenv("LOGGER_LEVEL")
	os.Unsetenv("HTTP_TIMEOUT")
	os.Unsetenv("SERVER_PORT")
}

func TestGetEnv(t *testing.T) {
	testKey := "TEST_KEY"
	expectedValue := "testValue"
	os.Setenv(testKey, expectedValue)

	assert.Equal(t, expectedValue, getEnv(testKey, "defaultValue"))
	assert.Equal(t, "defaultValue", getEnv("NON_EXISTENT_KEY", "defaultValue"))

	os.Unsetenv(testKey)
}

func TestGetEnvAsInt(t *testing.T) {
	testKey := "TEST_INT_KEY"
	expectedValue := 42
	os.Setenv(testKey, strconv.Itoa(expectedValue))

	assert.Equal(t, expectedValue, getEnvAsInt(testKey, 0))

	assert.Equal(t, 0, getEnvAsInt("NON_EXISTENT_KEY", 0))

	os.Unsetenv(testKey)
	assert.Equal(t, 0, getEnvAsInt(testKey, 0))
}
