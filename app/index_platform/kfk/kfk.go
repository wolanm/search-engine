package kfk

import (
	"github.com/pkg/errors"
	"github.com/wolanm/search-engine/consts"
	"github.com/wolanm/search-engine/kfk"
	"github.com/wolanm/search-engine/types"
)

func DocDataToKfk(fileInfo *types.FileInfo) error {
	data, _ := fileInfo.MarshalJSON()
	err := kfk.KafkaProducer(consts.KafkaIndexTopic, data)
	if err != nil {
		return errors.Wrap(err, "failed to produce document")
	}

	return nil
}
