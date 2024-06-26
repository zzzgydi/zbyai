package llm

import (
	"fmt"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"github.com/zzzgydi/zbyai/common/initializer"
)

var g3TokenEncoder *tiktoken.Tiktoken
var g4TokenEncoder *tiktoken.Tiktoken

func initTokenEncoders() error {
	gpt35TokenEncoder, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
	if err != nil {
		return fmt.Errorf("failed to get gpt-3.5-turbo token encoder: %w", err)
	}
	gpt4TokenEncoder, err := tiktoken.EncodingForModel("gpt-4")
	if err != nil {
		return fmt.Errorf("failed to get gpt-4 token encoder: %w", err)
	}

	g3TokenEncoder = gpt35TokenEncoder
	g4TokenEncoder = gpt4TokenEncoder
	return nil
}

func init() {
	initializer.Register("llm-token", initTokenEncoders)
}

func getTokenNum(tokenEncoder *tiktoken.Tiktoken, text string) int {
	return len(tokenEncoder.Encode(text, nil, nil))
}

// 粗略计算，不保真
func CountTokenMessages(model string, messages []openai.ChatCompletionMessage) int {
	// Reference:
	// https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	// https://github.com/pkoukk/tiktoken-go/issues/6
	//
	// Every message follows <|start|>{role/name}\n{content}<|end|>\n

	// TODO: 根据model选encoder
	encoder := g3TokenEncoder
	if strings.Contains(model, "gpt-4") {
		encoder = g4TokenEncoder
	}

	tokenNum := 0
	for _, message := range messages {
		tokenNum += 4
		tokenNum += getTokenNum(encoder, message.Content)
	}
	tokenNum += 3        // Every reply is primed with <|start|>assistant<|message|>
	return tokenNum + 36 // 随意写一个值，避免输入过长有偏差
}

// 粗略计算，不保真
func CountTokenText(model string, message string) int {
	// Reference:
	// https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	// https://github.com/pkoukk/tiktoken-go/issues/6
	//
	// Every message follows <|start|>{role/name}\n{content}<|end|>\n

	// TODO: 根据model选encoder
	encoder := g3TokenEncoder
	if strings.Contains(model, "gpt-4") {
		encoder = g4TokenEncoder
	}

	tokenNum := getTokenNum(encoder, message)
	tokenNum += 3        // Every reply is primed with <|start|>assistant<|message|>
	return tokenNum + 36 // 随意写一个值，避免输入过长有偏差
}
