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

type AppConfig struct {
	LLM *LLMConfig
	DB  *DbConfig
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

	return &AppConfig{
		llmConfig,
		dbConfig,
	}, nil
}
