package inverted_db

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/wolanm/search-engine/consts"
	"github.com/wolanm/search-engine/repository/boltdb"
	bolt "go.etcd.io/bbolt"
	"os"
)

type InvertedDB struct { // TODO: 后续做 mmap
	file   *os.File
	db     *bolt.DB
	offset int64
}

func NewInvertedDB(dbPath string) (*InvertedDB, error) {
	f, err := os.OpenFile(dbPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open %s failed", dbPath))
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get %s stat failed", dbPath))
	}

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open db failed"))
	}

	return &InvertedDB{
		file:   f,
		db:     db,
		offset: stat.Size(),
	}, nil
}

func (t *InvertedDB) StorageInverted(word string, values []byte) error {
	err := t.PutInverted([]byte(word), values)
	return errors.Wrap(err, "put inverted failed")
}

func (t *InvertedDB) PutInverted(key, value []byte) (err error) {
	err = boltdb.Put(t.db, consts.InvertedBucket, key, value)
	return errors.WithMessage(err, "put error")
}

func (t *InvertedDB) Close() error {
	err := t.file.Close()
	if err != nil {
		return err
	}

	err = t.db.Close()
	if err != nil {
		return err
	}

	return nil
}
