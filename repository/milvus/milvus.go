package milvus

import (
	"context"
	"github.com/milvus-io/milvus-proto/go-api/v2/milvuspb"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/pkg/errors"
)

type VectorDB struct {
	cli *milvusclient.Client
}

func (d *VectorDB) Insert(collection string, data []string) error {
	dataColumn := column.NewColumnString("data", data)
	_, err := d.cli.Insert(context.Background(), milvusclient.NewColumnBasedInsertOption(collection).WithColumns(dataColumn))
	if err != nil {
		return errors.Wrap(err, "insert failed")
	}

	return nil
}

func Search(collection string, query string) ([]int64, error) {

}
