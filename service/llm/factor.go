package llm

func NewLLM(model string, display string, temperature float32) *LLM {
	return &LLM{
		Model:       model,
		Display:     display,
		Temperature: temperature,
	}
}
