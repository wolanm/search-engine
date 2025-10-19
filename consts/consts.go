package consts

// etcd 相关
const (
	ServiceDomain = "etcd://"
	ServicePrefix = "/service"

	// etcd 服务名
	IndexPlatform = "index_platform"

	LeaseTTL = 10 // etcd 租约到期时间
)

// 服务名
const (
	GatewayService       = "gateway_service"
	IndexPlatformService = "index_platform_service"
	InvertedIndexService = "inverted_index_service"
	VectorIndexService   = "vector_index_service"
	SearchEngineService  = "search_engine_service"
)

const (
	DefaultKvListCapacity = 1e3 // 默认的 kvlist 容量
)

// mapreduce 相关
const (
	ConcurrentMapWorker = 3
)

// 数据库相关
const (
	InvertedDbCount = 5           // 倒排索引 boltdb 的分片数
	InvertedBucket  = "inverted"  // 倒排索引存储桶
	TrieTreeBucket  = "trie_tree" // trie 树存储桶
)

// kafka 相关
const (
	KafkaAssignorRoundRobin = "roundrobin"
	KafkaAssignorSticky     = "sticky"
	KafkaAssignorRange      = "range"

	// Topic
	KafkaIndexTopic = "search-engine-index-topic"

	// GroupID
	InvertedIndexGroupID = "kafka-inverted-index-group-id"
)
