package files

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/milosrs/go-hls-server/files/model"
)

const videoFolder = "videos"

type Service interface {
	AppendChunk(f *model.File) (int64, error)
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

func (s *ServiceImpl) AppendChunk(f *model.File) (int64, error) {
	if err := tryCreateVideoFolder(); err != nil {
		return 0, err
	}

	videoFile := filepath.Join(videoFolder, f.Name)
	file, err := os.OpenFile(videoFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileContent, err := f.Data.Open()
	if err != nil {
		return 0, err
	}
	defer fileContent.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, fileContent)
	if _, err := file.Write(buf.Bytes()); err != nil {
		return 0, err
	}

	info, err := os.Stat(videoFile)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
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
