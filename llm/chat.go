package llm

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/pkg/errors"
)

// ChatStream 流式输出
func ChatStream(ctx context.Context, history []*schema.Message) (*schema.StreamReader[*schema.Message], error) {
	messages := make([]*schema.Message, 0, len(history)+1)
	messages = append(messages, schema.SystemMessage(octoPrompt))
	messages = append(messages, history...)

	stream, err := tcm.Stream(ctx, messages)
	return stream, errors.WithStack(err)
}
