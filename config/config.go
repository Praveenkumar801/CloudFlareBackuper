package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CloudFlare CloudFlareConfig `yaml:"cloudflare"`
	Discord    DiscordConfig    `yaml:"discord"`
	Telegram   TelegramConfig   `yaml:"telegram"`
	Backup     BackupConfig     `yaml:"backup"`
}

type CloudFlareConfig struct {
	URI         string `yaml:"uri"`
	Bucket      string `yaml:"bucket"`
	AccessKeyID string `yaml:"access_key_id"`
	SecretKey   string `yaml:"secret_key"`
	AccountID   string `yaml:"account_id"`
}

type DiscordConfig struct {
	WebhookURL string `yaml:"webhook_url"`
}

type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   string `yaml:"chat_id"`
}

type BackupConfig struct {
	Schedule       string   `yaml:"schedule"`
	Folders        []string `yaml:"folders"`
	NamePrefix     string   `yaml:"name_prefix"`
	RetentionLimit int      `yaml:"retention_limit"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.CloudFlare.URI == "" {
		return fmt.Errorf("cloudflare.uri is required")
	}
	if c.CloudFlare.Bucket == "" {
		return fmt.Errorf("cloudflare.bucket is required")
	}
	if c.CloudFlare.AccessKeyID == "" {
		return fmt.Errorf("cloudflare.access_key_id is required")
	}
	if c.CloudFlare.SecretKey == "" {
		return fmt.Errorf("cloudflare.secret_key is required")
	}
	if c.CloudFlare.AccountID == "" {
		return fmt.Errorf("cloudflare.account_id is required")
	}
	// At least one notification method must be configured
	if c.Discord.WebhookURL == "" && c.Telegram.BotToken == "" {
		return fmt.Errorf("at least one notification method (discord or telegram) must be configured")
	}
	// Validate Telegram config if provided
	if c.Telegram.BotToken != "" && c.Telegram.ChatID == "" {
		return fmt.Errorf("telegram.chat_id is required when telegram.bot_token is provided")
	}
	if c.Backup.Schedule == "" {
		return fmt.Errorf("backup.schedule is required")
	}
	if len(c.Backup.Folders) == 0 {
		return fmt.Errorf("backup.folders must contain at least one folder")
	}
	if c.Backup.NamePrefix == "" {
		c.Backup.NamePrefix = "backup"
	}
	return nil
}
