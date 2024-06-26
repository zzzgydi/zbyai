package rag

import (
	"fmt"
	"sort"

	"github.com/zzzgydi/zbyai/common/utils"
	"github.com/zzzgydi/zbyai/model"
	"github.com/zzzgydi/zbyai/service/llm"
)

// 直接通过搜索结果进行简单的rerunk
func SimpleRerunk(results model.SearchList, maxToken int) string {
	if len(results) == 0 {
		return ""
	}

	// rerank
	chunks := []innerChunk{}

	maxLen := len(results)

	for idx, item := range results {
		// 如果没有爬取到结果，分数在 0-20
		// 如果爬取到结果，按token进行排序
		score := 0.0
		text := ""
		base := float64(maxLen - idx)
		weight := float64(2.0 * maxLen)

		if item.Page != "" && item.Token > 0 {
			// 超过比例进行扣分
			score = 2*base + weight
			if item.Token > maxToken {
				score -= float64(item.Token/maxToken) * 1.1133
				text = utils.SubstringMid(item.Page, maxToken*2/3, 1, 2) // 去掉两边的
			} else {
				score -= float64(item.Token/maxToken) * 1.5
				text = item.Page
			}
		} else {
			// 根据搜索结果的顺序进行排序
			score = 2 * base
			text = item.Snippet
		}

		chunks = append(chunks, innerChunk{text, score})
	}

	// 降序，所以这里用大于号
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Score > chunks[j].Score
	})

	// 生成结果
	content := ""

	for idx, item := range chunks {
		temp := content + fmt.Sprintf("## Doc %d\n%s\n\n", idx+1, item.Text)
		if (llm.CountTokenText("gpt-3.5-turbo", temp)) >= maxToken {
			break
		}
		content = temp
	}

	return content
}
