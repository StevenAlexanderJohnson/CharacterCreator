package internal

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyLLMUrl = errors.New("the llm url field is required in your config")
)

type LLMConfig struct {
	URL                    string `yaml:"url"`
	ApiKey                 string `yaml:"api-key"`
	AdditionalSystemPrompt string `yaml:"additional-system-prompt"`
}

func (l *LLMConfig) validate() error {
	if l.URL == "" {
		return ErrEmptyLLMUrl
	}

	return nil
}

type AppConfig struct {
	LLM LLMConfig `yaml:"llm"`
}

func (a *AppConfig) validate() error {
	if err := a.LLM.validate(); err != nil {
		return err
	}

	return nil
}

func ParseAppConfig(configPath string) (*AppConfig, error) {
	config := AppConfig{}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}
