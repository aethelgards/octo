package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aethelgards/octo/structs"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var configPath = ".octo/config.yaml"

var OctoConfig *structs.OctoConfig

func LoadConfig(_ context.Context) error {
	getwd, err := os.Getwd()
	if err != nil {
		return err
	}
	configFilePath := filepath.Join(getwd, configPath)
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return errors.WithStack(err)
	}
	OctoConfig = &structs.OctoConfig{}
	err = yaml.Unmarshal(file, OctoConfig)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := validateConfig(OctoConfig); err != nil {
		return err
	}
	return nil
}

func validateConfig(cfg *structs.OctoConfig) error {
	if cfg.LLMConfig.Model == "" {
		return errors.New("llm model is required")
	}
	if cfg.LLMConfig.BaseURL == "" {
		return errors.New("llm base_url is required")
	}
	if cfg.LLMConfig.APIKey == "" {
		return errors.New("llm api_key is required")
	}
	if cfg.LLMConfig.Timeout <= 0 {
		return errors.New("llm timeout must be positive")
	}

	if cfg.LogConfig.Level < 1 || cfg.LogConfig.Level > 4 {
		return errors.New("log level must be between 1 and 4")
	}
	if cfg.LogConfig.Format != "text" && cfg.LogConfig.Format != "json" {
		return errors.New("log format must be 'text' or 'json'")
	}
	return nil
}
