package thread

import (
	"errors"
	"fmt"

	"github.com/zzzgydi/zbyai/common/config"
	"github.com/zzzgydi/zbyai/common/initializer"
	"github.com/zzzgydi/zbyai/common/utils"
	"github.com/zzzgydi/zbyai/service/prompt"
)

var threadPM *prompt.PromptManager

var (
	chatModels    *utils.Chooser[config.ModelConfig]
	rewriteModels *utils.Chooser[config.ModelConfig]
)

func initThread() error {
	pm, err := prompt.NewPromptManager("v1")
	if err != nil {
		return err
	}
	threadPM = pm

	chatModels, err = transformModels(config.AppConf.ChatModels)
	if err != nil {
		return fmt.Errorf("transform chat models: %w", err)
	}

	rewriteModels, err = transformModels(config.AppConf.RewriteModels)
	if err != nil {
		return fmt.Errorf("transform rewrite models: %w", err)
	}

	return nil
}

func init() {
	initializer.Register("thread", initThread)
}

func transformModels(models []config.ModelConfig) (*utils.Chooser[config.ModelConfig], error) {
	if len(models) == 0 {
		return nil, errors.New("no models found")
	}

	modelChoice := make([]utils.Choice[config.ModelConfig], len(models))
	for i, m := range models {
		var weight uint
		if m.Weight <= 0 {
			weight = 1
		} else {
			weight = uint(m.Weight)
		}

		if m.Display == "" {
			m.Display = m.Model
		}

		modelChoice[i] = utils.Choice[config.ModelConfig]{Item: m, Weight: weight}
	}

	return utils.NewChooser(modelChoice), nil
}
