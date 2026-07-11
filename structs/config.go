package structs

type LLMConfig struct {
	Model   string `yaml:"model"`
	BaseURL string `yaml:"base_url"`
	APIKey  string `yaml:"api_key"`
	Timeout int    `yaml:"timeout"` // timeout, second
}

type OctoConfig struct {
	LLMConfig LLMConfig `yaml:"llm"`
	LogConfig LogConfig `yaml:"log"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level     int    `yaml:"level"`
	Format    string `yaml:"format"`
	LogDir    string `yaml:"log_dir"`
	AddSource bool   `yaml:"add_source"`
	Console   bool   `yaml:"console"`
}
