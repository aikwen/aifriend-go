package graph

import (
	"context"
	"fmt"
	"io"
	"log"

	einomodel "github.com/cloudwego/eino/components/model"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type Graph struct {
	model    einomodel.ToolCallingChatModel
	tools    map[string]einotool.InvokableTool
	maxSteps int
	debug    bool
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

func (g *Graph) SetDebug(debug bool) {
	g.debug = debug
}

func (g *Graph) logf(format string, args ...any) {
	if !g.debug {
		return
	}
	log.Printf("[chat-graph] "+format, args...)
}

func (g *Graph) StreamAndCollect(
	ctx context.Context,
	messages []*schema.Message,
	onChunk func(string),
) (*schema.Message, error) {
	// debug 日志
	g.logf("model input messages count=%d", len(messages))
	for i, msg := range messages {
		g.logf("input[%d] role=%q name=%q tool_call_id=%q content=%q tool_calls=%+v",
			i, msg.Role, msg.Name, msg.ToolCallID, msg.Content, msg.ToolCalls)
	}


	sr, err := g.model.Stream(ctx, messages)
	if err != nil {
		return nil, err
	}
	defer sr.Close()

	var chunks []*schema.Message
	chunkIndex := 0
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

		//
		g.logf("stream chunk[%d] role=%q content=%q tool_calls=%+v",
		chunkIndex, chunk.Role, chunk.Content, chunk.ToolCalls)
		g.logf("stream chunk[%d] response meta=%+v", chunkIndex, chunk.ResponseMeta)

		chunks = append(chunks, chunk)

		if onChunk != nil && chunk.Content != "" {
			onChunk(chunk.Content)
		}

		chunkIndex++
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
		// 日志
		g.logf("========== step %d start ==========", step+1)
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

		// 日志
		g.logf("step %d assistant appended: content=%q tool_calls=%+v",
			step+1, assistantMsg.Content, assistantMsg.ToolCalls)
		g.logf("step %d response meta: %+v", step+1, assistantMsg.ResponseMeta)

		// 没有 tool call，结束
		if len(assistantMsg.ToolCalls) == 0 {
			g.logf("step %d end: no tool calls, finish", step+1)
			return messages, nil
		}

		// 执行 tool calls
		for i, tc := range assistantMsg.ToolCalls {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			toolName := tc.Function.Name
			toolArgs := tc.Function.Arguments

			// 日志
			g.logf("step %d tool_call[%d]: id=%q name=%q args=%q",
				step+1, i, tc.ID, toolName, toolArgs)

			toolImpl, ok := g.tools[toolName]
			if !ok {
				return nil, fmt.Errorf("tool %q not found", toolName)
			}

			result, err := toolImpl.InvokableRun(ctx, toolArgs)
			if err != nil {
				return nil, fmt.Errorf("tool %q run failed: %w", toolName, err)
			}

			//
			g.logf("step %d tool_result[%d]: name=%q result=%q",
				step+1, i, toolName, result)

			toolMsg := &schema.Message{
				Role:       schema.Tool,
				Content:    result,
				ToolCallID: tc.ID,
				Name:       toolName,
			}

			messages = append(messages, toolMsg)

			g.logf("step %d tool message appended: tool_call_id=%q name=%q",
				step+1, tc.ID, toolName)
		}

		g.logf("========== step %d end ==========", step+1)
	}

	return nil, fmt.Errorf("max steps exceeded: %d", g.maxSteps)
}