package thread

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/zzzgydi/zbyai/common/utils"
	"github.com/zzzgydi/zbyai/service/llm"
)

type PeSearchJudge struct {
	UseSearch     bool   `json:"use_search,omitempty"`
	Model         string `json:"model,omitempty"`
	IsProgramming bool   `json:"is_programming,omitempty"`
	Language      string `json:"language,omitempty"`
	Rephrased     any    `json:"rephrased,omitempty"`
}

// 判断是否需要search
func peSearchJudge(query string, history []openai.ChatCompletionMessage, logger *slog.Logger) (*PeSearchJudge, error) {
	context := ""
	for _, item := range history {
		content := utils.RemoveMultiSpace(utils.Ellipsis(item.Content, 200))
		content = strings.TrimSpace(content)
		if item.Role == "user" {
			context += "<Q>" + content + "</Q>\n"
		} else {
			context += "<A>" + content + "</A>\n"
		}
	}

	query = utils.Ellipsis(query, 2000)

	generatePrompt := threadPM.GetPrompt("rewrite.txt")
	system := generatePrompt.Parse(nil)

	userQuery := fmt.Sprintf("<question>\n%s\n<Q>%s</Q>\n</question>", context, query)

	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: userQuery},
	}

	llmClient := LLMForJudge(0.52)

	output, err := llmClient.CompletionWithRetry(messages, 2)
	if err != nil {
		return nil, err
	}

	if logger != nil {
		// logger.Info("[Thread] pe search judge messages", "model", llmClient.Model, "messages", messages)
		logger.Info("[Thread] pe search judge", "model", llmClient.Model, "output", output)
	}

	res := &PeSearchJudge{}
	err = utils.TryParseJson(output, &res)
	if err != nil {
		return nil, fmt.Errorf("judge parse json error: %w", err)
	}

	res.Model = llmClient.Display

	return res, nil
}

func peMainRAG(query string, content string, history []openai.ChatCompletionMessage, logger *slog.Logger) (*llm.LLM, []openai.ChatCompletionMessage) {
	generatePrompt := threadPM.GetPrompt("generate.txt")

	if content == "" {
		content = "No context..."
	}
	system := generatePrompt.Parse(map[string]any{
		"context": content,
	})

	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: system},
	}
	if len(history) > 0 {
		messages = append(messages, history...)
	}
	messages = append(messages,
		openai.ChatCompletionMessage{Role: "user", Content: query})

	llmClient := LLMForMain(0.5)

	if logger != nil {
		usage := llm.CountTokenMessages(llmClient.Model, messages)
		// logger.Info("[Thread] main rag messages", "messages", messages)
		logger.Info("[Thread] main rag usage", "model", llmClient.Model, "pe_usage", usage)
	}

	return llmClient, messages
}
