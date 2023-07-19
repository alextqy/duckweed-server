package Lib

import (
	"bufio"
	"io/ioutil"
	"os"
)

func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func FileMake(filePath string) (bool, string) {
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func FileRemove(filePath string) (bool, string) {
	err := os.Remove(filePath)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func FileRename(filePath string, newName string) (bool, string) {
	err := os.Rename(filePath, newName)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

func Filespec(filePath string) (bool, os.FileInfo) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, nil
	}
	return true, fileInfo
}

func FileRead(filePath string) (bool, string) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		return false, err.Error()
	} else {
		contentByte, readErr := ioutil.ReadAll(f)
		if readErr != nil {
			return false, readErr.Error()
		}
		return true, string(contentByte)
	}
}

func FileWrite(filePath string, content string) (bool, string) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return false, err.Error()
	} else {
		_, writeErr := f.Write([]byte(content))
		if writeErr != nil {
			return false, err.Error()
		} else {
			return true, ""
		}
	}
}

func FileWriteAppend(filePath string, content string) (bool, string) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		return false, err.Error()
	} else {
		write := bufio.NewWriter(f)
		write.WriteString(content)
		write.Flush()
		return true, ""
	}
}

func DirMake(dirPath string) (bool, string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return false, err.Error()
	} else {
		return true, ""
	}
}

func DirCheck(dirPath string) (bool, string, []os.FileInfo) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "", files
}
