package model

type FileChunk struct {
	Name           string `json:"name" binding:"required"`
	Bytes          []byte `json:"bytes" binding:"required"`
	ChunkNumber    uint   `json:"chunkNumber" binding:"required"`
	NumberOfChunks int64  `json:"numberOfChunks" binding:"required"`
}

type InitialFileData struct {
	FileChunk
}
