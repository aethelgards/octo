package llm

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func AgentQuery(ctx context.Context, history []*schema.Message) *adk.AsyncIterator[*adk.AgentEvent] {
	return agentRunner.Run(ctx, history)
}
