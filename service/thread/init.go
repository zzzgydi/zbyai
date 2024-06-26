package thread

import (
	"github.com/zzzgydi/zbyai/common/initializer"
	"github.com/zzzgydi/zbyai/service/prompt"
)

var threadPM *prompt.PromptManager

func initThread() error {
	pm, err := prompt.NewPromptManager("v1")
	if err != nil {
		return err
	}
	threadPM = pm
	return nil
}

func init() {
	initializer.Register("thread", initThread)
}
