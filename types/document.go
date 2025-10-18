package types

type Document struct {
	DocId int64  `json:"doc_id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type FileInfo struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}
