package trie_db

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/wolanm/search-engine/consts"
	"github.com/wolanm/search-engine/repository/boltdb"
	"github.com/wolanm/search-engine/util/trie"
	bolt "go.etcd.io/bbolt"
	"os"
)

type TrieDB struct {
	file *os.File
	db   *bolt.DB
}

// NewTrieDB 初始化trie
func NewTrieDB(filePath string) (*TrieDB, error) { // TODO: 先都放在一个下面吧，后面再lb到多个文件
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open %s failed", filePath))
	}

	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open db failed"))
	}

	return &TrieDB{f, db}, nil
}

func (d *TrieDB) StorageDict(trieTree *trie.Trie) (err error) {
	trieByte, _ := trieTree.Root.Children.MarshalJSON()
	err = d.PutTrieTree([]byte(consts.TrieTreeBucket), trieByte)
	return errors.WithMessage(err, "putTrieTree error")
}

// GetTrieTreeInfo 获取 trie tree
func (d *TrieDB) GetTrieTreeInfo() (trieTree *trie.Trie, err error) {
	v, err := d.GetTrieTree([]byte(consts.TrieTreeBucket))
	if err != nil {
		err = errors.WithMessage(err, "getTrieTree error")
		return
	}

	trieTree = trie.NewTrie()
	err = trieTree.UnmarshalJSON(v)

	return trieTree, errors.Wrap(err, "failed to unmarshal data")
}

// PutTrieTree 存储
func (d *TrieDB) PutTrieTree(key, value []byte) (err error) {
	err = boltdb.Put(d.db, consts.TrieTreeBucket, key, value)
	return errors.WithMessage(err, "put error")
}

// GetTrieTree 通过term获取value
func (d *TrieDB) GetTrieTree(key []byte) (value []byte, err error) {
	value, err = boltdb.Get(d.db, consts.TrieTreeBucket, key)
	if err != nil {
		err = errors.WithMessage(err, "get error")
	}
	return
}

// Close 关闭db
func (d *TrieDB) Close() error {
	return errors.Wrap(d.db.Close(), "failed to close")
}
