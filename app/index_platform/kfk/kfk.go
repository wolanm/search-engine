package kfk

import (
	"github.com/pkg/errors"
	"github.com/wolanm/search-engine/consts"
	"github.com/wolanm/search-engine/kfk"
	"github.com/wolanm/search-engine/types"
)

func DocDataToKfk(doc *types.Document) error {
	docBytes, _ := doc.MarshalJSON()
	err := kfk.KafkaProducer(consts.KafkaIndexTopic, docBytes)
	if err != nil {
		return errors.Wrap(err, "failed to produce document")
	}

	return nil
}
