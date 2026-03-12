package graph

import (
	"context"
	"fmt"
	"io"

	einomodel "github.com/cloudwego/eino/components/model"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type Graph struct {
	model    einomodel.ToolCallingChatModel
	tools    map[string]einotool.InvokableTool
	maxSteps int
}

func NewGraph(
	model einomodel.ToolCallingChatModel,
	tools map[string]einotool.InvokableTool,
) *Graph {
	return &Graph{
		model:    model,
		tools:    tools,
		maxSteps: 10,
	}
}

func (g *Graph) SetMaxSteps(n int) {
	if n > 0 {
		g.maxSteps = n
	}
}

func (g *Graph) StreamAndCollect(
	ctx context.Context,
	messages []*schema.Message,
	onChunk func(string),
) (*schema.Message, error) {
	sr, err := g.model.Stream(ctx, messages)
	if err != nil {
		return nil, err
	}
	defer sr.Close()

	var chunks []*schema.Message

	for {
		chunk, err := sr.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if chunk == nil {
			continue
		}

		chunks = append(chunks, chunk)

		if onChunk != nil && chunk.Content != "" {
			onChunk(chunk.Content)
		}
	}

	return schema.ConcatMessages(chunks)
}

func (g *Graph) Run(
	ctx context.Context,
	messages []*schema.Message,
	onChunk func(string),
) ([]*schema.Message, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("messages is empty")
	}

	for step := 0; step < g.maxSteps; step++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		assistantMsg, err := g.StreamAndCollect(ctx, messages, onChunk)
		if err != nil {
			return nil, fmt.Errorf("model stream failed at step %d: %w", step+1, err)
		}
		if assistantMsg == nil {
			return nil, fmt.Errorf("assistant message is nil at step %d", step+1)
		}

		messages = append(messages, assistantMsg)

		// 没有 tool call，结束
		if len(assistantMsg.ToolCalls) == 0 {
			return messages, nil
		}

		// 执行 tool calls
		for _, tc := range assistantMsg.ToolCalls {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			toolName := tc.Function.Name
			toolArgs := tc.Function.Arguments

			toolImpl, ok := g.tools[toolName]
			if !ok {
				return nil, fmt.Errorf("tool %q not found", toolName)
			}

			result, err := toolImpl.InvokableRun(ctx, toolArgs)
			if err != nil {
				return nil, fmt.Errorf("tool %q run failed: %w", toolName, err)
			}

			toolMsg := &schema.Message{
				Role:       schema.Tool,
				Content:    result,
				ToolCallID: tc.ID,
				Name:       toolName,
			}

			messages = append(messages, toolMsg)
		}
	}

	return nil, fmt.Errorf("max steps exceeded: %d", g.maxSteps)
}