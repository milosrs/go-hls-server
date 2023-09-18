package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/milosrs/go-hls-server/files/model"
)

const videoFolder = "videos"

type Service interface {
	CreateFile(f *model.InitialFileData) (float64, error)
	AppendChunk(f *model.FileChunk) (float64, error)
	Remove(name string) error
	Start()
	Stop()
}

type createFileReq struct {
	data *model.InitialFileData
	resp chan error
}

type appendFileReq struct {
	data *model.FileChunk
	resp chan error
}

type ServiceImpl struct {
	createFile chan createFileReq
	appendFile chan appendFileReq
	stop       chan struct{}
}

func NewService() Service {
	return &ServiceImpl{
		createFile: make(chan createFileReq),
		appendFile: make(chan appendFileReq),
		stop:       make(chan struct{}),
	}
}

func tryCreateVideoFolder() error {
	if _, err := os.Stat(videoFolder); os.IsNotExist(err) {
		return os.Mkdir(videoFolder, 0777)
	}

	return nil
}

func doesFileExist(name string) bool {
	_, err := os.Stat(videoFolder + "/" + name)

	if err != nil {
		return false
	}

	return true
}

func writeBytesToFile(f *os.File, bytes []byte) error {
	_, err := f.Write(bytes)
	if err != nil {
		return err
	}

	_, err = f.Stat()
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) Start() {
	for {
		select {
		case req := <-s.createFile:
			_, err := s.CreateFile(req.data)
			req.resp <- err
		case req := <-s.appendFile:
			_, err := s.AppendChunk(req.data)
			req.resp <- err
		case <-s.stop:
			s.Stop()
			return
		}
	}
}

func (s *ServiceImpl) Stop() {
	close(s.appendFile)
	close(s.createFile)
	close(s.stop)
}

func (s *ServiceImpl) CreateFile(f *model.InitialFileData) (float64, error) {
	if err := tryCreateVideoFolder(); err != nil {
		return 0, err
	}

	fileExists := doesFileExist(f.Name)
	if fileExists {
		info, err := os.Stat(videoFolder + "/" + f.Name)
		if err != nil {
			return 0, err
		}

		return float64(info.Size()), err
	}

	videoFile := filepath.Join(videoFolder, f.Name)
	file, err := os.OpenFile(videoFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0604)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	err = writeBytesToFile(file, f.Bytes)
	if err != nil {
		return 0, err
	}

	return float64(int64(f.ChunkNumber) / f.NumberOfChunks), nil
}

func (*ServiceImpl) AppendChunk(f *model.FileChunk) (float64, error) {
	videoFile := filepath.Join(videoFolder, f.Name)
	file, err := os.OpenFile(videoFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	err = writeBytesToFile(file, f.Bytes)
	if err != nil {
		return 0, err
	}

	return float64(int64(f.ChunkNumber) / f.NumberOfChunks), nil
}

func (s *ServiceImpl) Remove(name string) error {
	videoFile := filepath.Join(videoFolder, name)

	info, err := os.Stat(videoFile)
	if err != nil {
		return err
	}
	fmt.Println(info)

	return os.Remove(videoFile)
}
