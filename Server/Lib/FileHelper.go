package Lib

import (
	"bufio"
	"io/ioutil"
	"os"
)

func FileExist(FilePath string) bool {
	_, Err := os.Stat(FilePath)
	return !os.IsNotExist(Err)
}

func FileMake(FilePath string) (bool, string) {
	F, Err := os.Create(FilePath)
	defer F.Close()
	if Err != nil {
		return false, Err.Error()
	}
	return true, ""
}

func FileRemove(FilePath string) (bool, string) {
	Err := os.Remove(FilePath)
	if Err != nil {
		return false, Err.Error()
	}
	return true, ""
}

func FileRename(FilePath string, NewName string) (bool, string) {
	Err := os.Rename(FilePath, NewName)
	if Err != nil {
		return false, Err.Error()
	}
	return true, ""
}

func Filespec(FilePath string) (bool, os.FileInfo) {
	FileInfo, Err := os.Stat(FilePath)
	if Err != nil {
		return false, nil
	}
	return true, FileInfo
}

func FileRead(FilePath string) (bool, string) {
	F, Err := os.OpenFile(FilePath, os.O_RDONLY, 0600)
	defer F.Close()
	if Err != nil {
		return false, Err.Error()
	} else {
		ContentByte, ReadErr := ioutil.ReadAll(F)
		if ReadErr != nil {
			return false, ReadErr.Error()
		}
		return true, string(ContentByte)
	}
}

func FileWrite(FilePath string, Content string) (bool, string) {
	F, Err := os.OpenFile(FilePath, os.O_WRONLY|os.O_TRUNC, 0600)
	defer F.Close()
	if Err != nil {
		return false, Err.Error()
	} else {
		_, WriteErr := F.Write([]byte(Content))
		if WriteErr != nil {
			return false, Err.Error()
		} else {
			return true, ""
		}
	}
}

func FileWriteAppend(FilePath string, Content string) (bool, string) {
	F, Err := os.OpenFile(FilePath, os.O_WRONLY|os.O_APPEND, 0666)
	defer F.Close()
	if Err != nil {
		return false, Err.Error()
	} else {
		Write := bufio.NewWriter(F)
		Write.WriteString(Content)
		Write.Flush()
		return true, ""
	}
}

func DirMake(DirPath string) (bool, string) {
	Err := os.MkdirAll(DirPath, os.ModePerm)
	if Err != nil {
		return false, Err.Error()
	} else {
		return true, ""
	}
}

func DirCheck(DirPath string) (bool, string, []os.FileInfo) {
	Files, Err := ioutil.ReadDir(DirPath)
	if Err != nil {
		return false, Err.Error(), nil
	}
	return true, "", Files
}
