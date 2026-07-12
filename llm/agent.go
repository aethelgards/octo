package llm

import (
	"context"

	"github.com/cloudwego/eino/adk"
)

var agentRunner *adk.Runner

func InitAgent(ctx context.Context) error {
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "octo",
		Description: "A helpful assistant that can answer questions and help with tasks.",
		Instruction: octoPrompt,
		Model:       tcm,
	})
	if err != nil {
		return err
	}

	agentRunner = adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	return nil
}
