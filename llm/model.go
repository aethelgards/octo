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
	var chatModelConfig = &openai.ChatModelConfig{
		APIKey:  llmConfig.APIKey,
		Timeout: time.Second * time.Duration(llmConfig.Timeout),
		BaseURL: llmConfig.BaseURL,
		Model:   llmConfig.Model,
	}
	tcm, err = openai.NewChatModel(ctx, chatModelConfig)
	return err
}
