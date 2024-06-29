package thread

import (
	"github.com/zzzgydi/zbyai/service/llm"
)

func LLMForMain(temperature float32) *llm.LLM {
	t := chatModels.Pick()
	if t == nil {
		return nil
	}
	model := *t
	return llm.NewLLM(model.Model, model.Display, temperature)
}

func LLMForJudge(temperature float32) *llm.LLM {
	t := rewriteModels.Pick()
	if t == nil {
		return nil
	}
	model := *t
	return llm.NewLLM(model.Model, model.Display, temperature)
}
