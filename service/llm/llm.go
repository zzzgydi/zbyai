package llm

import (
	"context"
	"errors"
	"log/slog"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func (l *LLM) Completion(messages []openai.ChatCompletionMessage) (string, error) {
	client := openai.NewClientWithConfig(*gptConfig)

	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:       l.Model,
		Messages:    messages,
		Stream:      false,
		Temperature: l.Temperature,
	})
	if err != nil {
		return "", err
	}

	slog.Info("LLM usage", "model", l.Model, "usage", resp.Usage)

	if len(resp.Choices) == 0 {
		slog.Error("LLM completion no choices", "response", resp)
		return "", errors.New("response has no choices")
	}

	return resp.Choices[0].Message.Content, nil
}

func (l *LLM) CompletionWithRetry(messages []openai.ChatCompletionMessage, retry int) (string, error) {
	var retErr error
	for i := 0; i < retry; i++ {
		res, err := l.Completion(messages)
		if err == nil {
			return res, nil
		}

		slog.Error("Completion With Retry", "retry", i, "error", err)
		retErr = err
		time.Sleep(time.Second * 2)
	}
	return "", retErr
}

func (l *LLM) StreamCompletion(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	client := openai.NewClientWithConfig(*gptConfig)

	req := openai.ChatCompletionRequest{
		Model:       l.Model,
		Messages:    messages,
		Stream:      true,
		Temperature: l.Temperature,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	return stream, err
}
