package prompt

import (
	"log/slog"
	"os"
	"path/filepath"

	tmp "github.com/zzzgydi/templater"
	"github.com/zzzgydi/zbyai/common/initializer"
)

var (
	promptCache map[string]map[string]*tmp.Templater
)

func InitPrompt() error {
	promptDir := "prompt"
	promptCache = make(map[string]map[string]*tmp.Templater)

	// 读取prompt目录下所有的子目录
	err := filepath.Walk(promptDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			promptCache[info.Name()], _ = readPromptDir(path)
		}
		return nil
	})

	return err
}

func readPromptDir(dirPath string) (map[string]*tmp.Templater, error) {
	promptMap := make(map[string]*tmp.Templater)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				slog.Error("read file error", "file", path, "error", err)
			} else {
				promptMap[info.Name()] = tmp.NewTemplater(string(content))
			}
		}
		return nil
	})
	return promptMap, err
}

func init() {
	initializer.Register("prompt", InitPrompt)
}
