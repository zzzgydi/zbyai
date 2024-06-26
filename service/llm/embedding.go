package llm

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

func TextEmbedding(ctx context.Context, inputs []string) ([]openai.Embedding, error) {
	client := openai.NewClientWithConfig(*gptConfig)
	res, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequestStrings{
		Input: inputs,
		Model: openai.SmallEmbedding3,
	})
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func QueryEmbedding(ctx context.Context, inputs []string) ([]openai.Embedding, error) {
	temps := make([]string, len(inputs))
	for i, input := range inputs {
		temps[i] = "Query: " + input + "?"
	}
	return TextEmbedding(ctx, temps)
}

func AnswerEmbedding(ctx context.Context, inputs []string) ([]openai.Embedding, error) {
	temps := make([]string, len(inputs))
	for i, input := range inputs {
		temps[i] = "Answer: " + input + "."
	}
	return TextEmbedding(ctx, temps)
}
