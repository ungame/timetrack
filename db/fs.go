package db

import (
	"github.com/ungame/timetrack/ioext"
	"log"
	"os"
	"path"
	"runtime"
)

const (
	defaultDir  = ".data"
	defaultFile = "data.sqlite"
)

type FileStorage interface {
	Create() string
}

type defaultFileStorage struct{}

func DefaultFileStorage() FileStorage {
	return &defaultFileStorage{}
}

func (d *defaultFileStorage) Create() string {
	var (
		currentPath  = getCurrentPath()
		fullDirPath  = currentPath + "/" + defaultDir
		fullFilePath = fullDirPath + "/" + defaultFile
	)

	if _, err := os.Stat(fullDirPath); os.IsNotExist(err) {
		if err = os.Mkdir(fullDirPath, os.ModePerm); err != nil {
			log.Panicln("error on create storage:", err)
		}
	}

	if _, err := os.Stat(fullFilePath); os.IsNotExist(err) {
		if file, err := os.Create(fullFilePath); err == nil {
			ioext.Close(file)
		} else {
			log.Println("error on create file:", err.Error())
		}
	}

	return fullFilePath
}

func getCurrentPath() string {
	_, file, _, _ := runtime.Caller(0)
	return path.Dir(file)
}
