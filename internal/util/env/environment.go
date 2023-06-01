package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type IEnvironment interface {
	Init()
	Get(key string) string
	Set(key string, value string) error
	GetHostname() (string, error)
}

type Environment struct{}

// New
// Returns new Environment.
func New() IEnvironment {
	return &Environment{}
}

func (e *Environment) Init() {
	// Load env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err.Error())
		panic("Panicked while loading environment.")
	}

	appEnv := os.Getenv("APP_ENVIRONMENT")
	appName := os.Getenv("APP_NAME")
	if len(appEnv) < 1 {
		panic("APP_ENVIRONMENT variable is not set.")
	}
	if len(appName) < 1 {
		panic("APP_NAME variable is not set.")
	}
}

func (e *Environment) Get(key string) string {
	return os.Getenv(key)
}

func (e *Environment) Set(key string, value string) error {
	return os.Setenv(key, value)
}

func (e *Environment) GetHostname() (string, error) {
	return os.Hostname()
}
