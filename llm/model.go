package llm

import (
	"context"
	"time"

	"github.com/aethelgards/octo/structs"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

var tcm model.ToolCallingChatModel

func InitModel(ctx context.Context, llmConfig structs.LLMConfig) (err error) {
	var extMap = make(map[string]any)
	if llmConfig.Thinking.Enabled {
		extMap["thinking"] = map[string]string{"type": "enabled"}
		if llmConfig.Thinking.ReasoningEffort != "" {
			extMap["reasoning_effort"] = llmConfig.Thinking.ReasoningEffort
		}
	}
	var chatModelConfig = &openai.ChatModelConfig{
		APIKey:      llmConfig.APIKey,
		Timeout:     time.Second * time.Duration(llmConfig.Timeout),
		BaseURL:     llmConfig.BaseURL,
		Model:       llmConfig.Model,
		ExtraFields: extMap,
	}
	tcm, err = openai.NewChatModel(ctx, chatModelConfig)
	return err
}
