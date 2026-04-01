package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	LLM struct {
		Provider string `mapstructure:"provider"`
		OpenAI   struct {
			APIKey  string `mapstructure:"api_key"`
			Model   string `mapstructure:"model"`
			BaseURL string `mapstructure:"base_url"`
		} `mapstructure:"openai"`
		Anthropic struct {
			APIKey  string `mapstructure:"api_key"`
			Model   string `mapstructure:"model"`
			BaseURL string `mapstructure:"base_url"`
		} `mapstructure:"anthropic"`
		Ollama struct {
			BaseURL string `mapstructure:"base_url"`
			Model   string `mapstructure:"model"`
		} `mapstructure:"ollama"`
	} `mapstructure:"llm"`
	Security struct {
		FilterPatterns []FilterPattern `mapstructure:"filter_patterns"`
	} `mapstructure:"security"`
	UI struct {
		Theme   string `mapstructure:"theme"`
		Hotkeys struct {
			CopyLastCommand string `mapstructure:"copy_last_command"`
			FocusApp        string `mapstructure:"focus_app"`
		} `mapstructure:"hotkeys"`
	} `mapstructure:"ui"`
}

// FilterPattern defines a regex pattern for sensitive data.
type FilterPattern struct {
	Name    string `mapstructure:"name"`
	Pattern string `mapstructure:"pattern"`
}

var (
	globalConfig *Config
	configPath   string
)

// Init loads configuration from file and environment variables.
func Init(configFile string) error {
	configPath = configFile
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("llm.provider", "openai")
	viper.SetDefault("llm.openai.model", "gpt-4")
	viper.SetDefault("llm.openai.base_url", "https://api.openai.com/v1")
	viper.SetDefault("llm.anthropic.model", "claude-3-haiku-20240307")
	viper.SetDefault("llm.anthropic.base_url", "https://api.anthropic.com")
	viper.SetDefault("llm.ollama.base_url", "http://localhost:11434")
	viper.SetDefault("llm.ollama.model", "llama3")
	viper.SetDefault("ui.theme", "system")
	viper.SetDefault("ui.hotkeys.copy_last_command", "Ctrl+Shift+C")
	viper.SetDefault("ui.hotkeys.focus_app", "Ctrl+Shift+P")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			// Create directory if it doesn't exist
			dir := filepath.Dir(configFile)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			if err := viper.WriteConfigAs(configFile); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Bind environment variables
	viper.SetEnvPrefix("PAIRADMIN")
	viper.AutomaticEnv()

	// Unmarshal
	globalConfig = &Config{}
	if err := viper.Unmarshal(globalConfig); err != nil {
		return err
	}

	return nil
}

// Get returns the global configuration.
func Get() *Config {
	return globalConfig
}

// Save writes the current configuration to disk.
func Save() error {
	viper.Set("llm", globalConfig.LLM)
	viper.Set("security", globalConfig.Security)
	viper.Set("ui", globalConfig.UI)
	return viper.WriteConfig()
}