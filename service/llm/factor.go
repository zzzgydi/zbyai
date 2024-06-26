package llm

import "github.com/sashabaranov/go-openai"

func NewLLM(model string, temperature float32) *LLM {
	return &LLM{
		Model:       model,
		Temperature: temperature,
	}
}

func NewLLMGPT35(temperature float32) *LLM {
	return &LLM{
		Model:       openai.GPT3Dot5Turbo,
		Temperature: temperature,
	}
}

func NewLLMGPT4(temperature float32) *LLM {
	return &LLM{
		Model:       openai.GPT4Turbo1106,
		Temperature: temperature,
	}
}

func NewLLMGroq(temperature float32) *LLM {
	return &LLM{
		Model:       "mixtral-8x7b-32768",
		Temperature: temperature,
	}
}

func NewLLMGeminiPro(temperature float32) *LLM {
	return &LLM{
		Model:       "gemini-1.0-pro-001",
		Temperature: temperature,
	}
}

func NewLLMGemini15Pro(temperature float32) *LLM {
	return &LLM{
		Model:       "gemini-1.5-pro-latest",
		Temperature: temperature,
	}
}

func NewLLMClaudeHaiku(temp float32) *LLM {
	return &LLM{
		Model:       "anthropic/claude-3-haiku",
		Temperature: temp,
	}
}
