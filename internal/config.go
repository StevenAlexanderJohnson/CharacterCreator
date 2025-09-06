package internal

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrEmptyLLMUrl = errors.New("the llm url field is required in your config")
)

type LLMConfig struct {
	URL                    string
	ApiKey                 string
	AdditionalSystemPrompt string
}

func LoadLLMConfigEnv() (*LLMConfig, error) {
	url := os.Getenv("LLM_URL")
	if url == "" {
		return nil, fmt.Errorf("no llm db url was provided")
	}
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("no llm db api key was provided")
	}
	systemPrompt := os.Getenv("LLM_SYSTEM_PROMPT")

	return &LLMConfig{
		URL:                    url,
		ApiKey:                 apiKey,
		AdditionalSystemPrompt: systemPrompt,
	}, nil
}

type DbConfig struct {
	SqlitePath string
}

func LoadDbConfigEnv() (*DbConfig, error) {
	dbFilePath := os.Getenv("DB_FILE_PATH")
	if dbFilePath == "" {
		return nil, fmt.Errorf("no db file path provided")
	}
	return &DbConfig{
		SqlitePath: dbFilePath,
	}, nil
}

type AuthServiceConfig struct {
	URL         string
	ServiceName string
}

func LoadAuthServiceConfigEnv() (*AuthServiceConfig, error) {
	authServiceURL := os.Getenv("AUTH_URL")
	if authServiceURL == "" {
		return nil, fmt.Errorf("no auth service url provided")
	}
	serviceName := os.Getenv("AUTH_SERVICE_NAME")
	if serviceName == "" {
		return nil, fmt.Errorf("no service name was provided")
	}

	return &AuthServiceConfig{
		URL:         authServiceURL,
		ServiceName: serviceName,
	}, nil
}

type AppConfig struct {
	LLM               *LLMConfig
	DB                *DbConfig
	AuthServiceConfig *AuthServiceConfig
}

func ParseAppConfig() (*AppConfig, error) {

	llmConfig, err := LoadLLMConfigEnv()
	if err != nil {
		return nil, err
	}

	dbConfig, err := LoadDbConfigEnv()
	if err != nil {
		return nil, err
	}

	authServiceConfig, err := LoadAuthServiceConfigEnv()
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		llmConfig,
		dbConfig,
		authServiceConfig,
	}, nil
}
