package consts

const (
	// etcd 相关
	ServiceDomain = "etcd://"
	ServicePrefix = "/service"

	// etcd 服务名
	IndexPlatform = "index_platform"

	LeaseTTL = 10 // etcd 租约到期时间

	DefaultKvListCapacity = 1e3 // 默认的 kvlist 容量

	// 数据库相关
	InvertedDbCount = 5           // 倒排索引 boltdb 的分片数
	InvertedBucket  = "inverted"  // 倒排索引存储桶
	TrieTreeBucket  = "trie_tree" // trie 树存储桶
)
