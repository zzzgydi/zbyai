package thread

import (
	"github.com/zzzgydi/zbyai/common/utils"
	"github.com/zzzgydi/zbyai/service/llm"
)

var (
	mainModels  *utils.Chooser[string]
	judgeModels *utils.Chooser[string]
)

func init() {
	mainModels = utils.NewChooser([]utils.Choice[string]{
		{Item: "gemini-pro", Weight: 6},
		{Item: "cohere/command-r", Weight: 10},
		{Item: "anthropic/claude-3-haiku", Weight: 10},
		{Item: "gpt-3.5-turbo", Weight: 15},
		{Item: "meta-llama/llama-3-8b-instruct", Weight: 20},
		{Item: "deepseek/deepseek-chat", Weight: 25},
	})
	judgeModels = utils.NewChooser([]utils.Choice[string]{
		{Item: "gpt-3.5-turbo-0125", Weight: 50},
		{Item: "deepseek/deepseek-chat", Weight: 10},
		{Item: "anthropic/claude-3-haiku", Weight: 5},
		{Item: "gemini-pro", Weight: 5},
		{Item: "gpt-4-1106-preview", Weight: 1},
		{Item: "meta-llama/llama-3-8b-instruct", Weight: 10},
	})
}

func LLMForMain(temperature float32) *llm.LLM {
	model := mainModels.Pick()
	return llm.NewLLM(*model, temperature)
}

func LLMForJudge(temperature float32) *llm.LLM {
	model := judgeModels.Pick()
	return llm.NewLLM(*model, temperature)
}
