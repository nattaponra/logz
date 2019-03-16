package logz

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/lestrrat-go/strftime"
)

type Logz struct {
	fileName        string
	pathFile        *strftime.Strftime
	locker          sync.Locker
	file            *os.File
	symlink         *strftime.Strftime
	dateTimePattern string
}

func NewLogz(fileName, dateTimePattern string) *Logz {

	pathFile, err := strftime.New(fileName + "." + dateTimePattern)
	if err != nil {
		fmt.Println(err.Error())
	}

	return &Logz{
		pathFile:        pathFile,
		fileName:        fileName,
		dateTimePattern: dateTimePattern,
		locker:          new(sync.Mutex),
		symlink:         nil,
	}
}

func (logz *Logz) Write(b []byte) error {

	logz.locker.Lock()
	defer logz.locker.Unlock()

	go func(fp *os.File) {
		if fp == nil {
			return
		}
		fp.Close()
	}(logz.file)

	filePath := logz.pathFile.FormatString(time.Now())
	fmt.Println(filePath)
	//Create directory.
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	fp, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error(err.Error())
	}

	//defer fp.Close()

	logz.createSymlink(time.Now(), filePath)
	_, err = fp.Write(b)

	return err
}

func (logz *Logz) createSymlink(t time.Time, path string) {

	symlinkPath, err := strftime.New(logz.fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	symlink := symlinkPath.FormatString(t)
	fmt.Println(symlink)
	if symlink == path {
		return
	}

	if _, err := os.Stat(symlink); err == nil {
		if err := os.Remove(symlink); err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if err := os.Symlink(path, symlink); err != nil {
		fmt.Println(err.Error())
		return
	}
}
