package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/milosrs/go-hls-server/files/model"
)

const videoFolder = "videos"

type Service interface {
	CreateFile(f *model.InitialFileData) (int64, error)
	AppendChunk(f *model.FileChunk) (int64, error)
	Remove(name string) error
}

type ServiceImpl struct {
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

func writeBytesToFile(f *os.File, bytes []byte) (int64, error) {
	_, err := f.Write(bytes)
	if err != nil {
		return 0, err
	}

	info, err := f.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func (s *ServiceImpl) CreateFile(f *model.InitialFileData) (int64, error) {
	if err := tryCreateVideoFolder(); err != nil {
		return 0, err
	}

	fileExists := doesFileExist(f.Name)
	if fileExists {
		info, err := os.Stat(videoFolder + "/" + f.Name)
		if err != nil {
			return 0, err
		}

		return info.Size(), err
	}

	videoFile := filepath.Join(videoFolder, f.Name)
	file, err := os.OpenFile(videoFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0604)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return writeBytesToFile(file, f.Bytes)
}

func (s *ServiceImpl) AppendChunk(f *model.FileChunk) (int64, error) {
	videoFile := filepath.Join(videoFolder, f.Name)
	file, err := os.OpenFile(videoFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return writeBytesToFile(file, f.Bytes)
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
