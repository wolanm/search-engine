package input_data

import (
	"github.com/spf13/cast"
	"github.com/wolanm/search-engine/types"
	"strings"
)

func Doc2Struct(docStr string) (*types.Document, error) {
	docStr = strings.Replace(docStr, "\"", "", -1)
	d := strings.Split(docStr, ",")
	something2Str := make([]string, 0)

	for i := 2; i < 5; i++ {
		if len(d) > i && d[i] != "" {
			something2Str = append(something2Str, d[i])
		}
	}

	doc := &types.Document{
		DocId: cast.ToInt64(d[0]),
		Title: d[1],
		Body:  strings.Join(something2Str, ""),
	}

	return doc, nil
}
