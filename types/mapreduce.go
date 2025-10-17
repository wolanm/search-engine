package types

type Tokenization struct {
	Token string // 词条
	// Position int64  // 词条在文本的位置 // TODO 后面再补上
	// Offset   int64  // 偏移量
	DocId int64
}
