package structs

type ThinkingConfig struct {
	Enabled         bool   `yaml:"enabled"`
	ReasoningEffort string `yaml:"reasoning_effort"` // "low" | "medium" | "high"
	Show            bool   `yaml:"show"`
}

type LLMConfig struct {
	Model    string         `yaml:"model"`
	BaseURL  string         `yaml:"base_url"`
	APIKey   string         `yaml:"api_key"`
	Timeout  int            `yaml:"timeout"` // timeout, second
	Thinking ThinkingConfig `yaml:"thinking"`
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
