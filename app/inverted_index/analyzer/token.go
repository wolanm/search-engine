package analyzer

import (
	"github.com/go-ego/gse"
	"github.com/wolanm/search-engine/types"
	"strings"
)

var (
	GlobalSega gse.Segmenter
)

// InitSeg 分词器初始化
func InitSeg() {
	newGse, _ := gse.New()
	GlobalSega = newGse
}

// GseCutForBuildIndex 分词 IK for building index
func GseCutForBuildIndex(docId int64, content string) ([]*types.Tokenization, error) {
	content = ignoredChar(content)
	c := GlobalSega.CutSearch(content)
	token := make([]*types.Tokenization, 0)
	for _, v := range c {
		token = append(token, &types.Tokenization{
			Token: v,
			DocId: docId,
		})
	}

	return token, nil
}

func ignoredChar(str string) string {
	for _, c := range str {
		switch c {
		case '\f', '\n', '\r', '\t', '\v', '!', '"', '#', '$', '%', '&',
			'\'', '(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>',
			'?', '@', '[', '\\', '【', '】', ']', '“', '”', '「', '」', '★', '^', '·', '_', '`', '{', '|', '}', '~', '《', '》', '：',
			'（', '）', 0x3000, 0x3001, 0x3002, 0xFF01, 0xFF0C, 0xFF1B, 0xFF1F:
			str = strings.ReplaceAll(str, string(c), "")
		}
	}
	return str
}
