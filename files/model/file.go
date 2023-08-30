package model

type FileChunk struct {
	Name  string `json:"name" binding:"required"`
	Bytes []byte `json:"bytes" binding:"required"`
}

type InitialFileData struct {
	FileChunk
	NumberOfChunks int64 `json:"numberOfChunks" binding:"required"`
}
