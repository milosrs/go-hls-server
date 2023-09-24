package model

type FileChunk struct {
	Bytes []byte `json:"bytes" binding:"required"`
	FileDescription
}

type InitialFileData struct {
	FileChunk
}

type FileDescription struct {
	Name           string `json:"name" binding:"required"`
	ChunkNumber    uint   `json:"chunkNumber" binding:"required"`
	NumberOfChunks int64  `json:"numberOfChunks" binding:"required"`
}
