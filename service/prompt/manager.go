package prompt

import (
	"fmt"

	tmp "github.com/zzzgydi/templater"
)

type PromptManager struct {
	typ string
}

func NewPromptManager(typ string) (*PromptManager, error) {
	if _, ok := promptCache[typ]; !ok {
		return nil, fmt.Errorf("prompt type not found")
	}
	return &PromptManager{
		typ: typ,
	}, nil
}

func (p *PromptManager) GetPrompt(promptType string) *tmp.Templater {
	return promptCache[p.typ][promptType]
}
