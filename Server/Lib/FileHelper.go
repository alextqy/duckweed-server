package lib

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// 文件检查
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// 新建文件
func FileMake(filePath string) (bool, string) {
	f, err := os.Create(filePath)

	defer func(f io.Closer) {
		if err := f.Close(); err != nil {
			fmt.Printf("defer close file err: %v", err.Error())
		}
	}(f)

	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

// 文件删除
func FileRemove(filePath string) (bool, string) {
	err := os.Remove(filePath)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

// 文件重命名
func FileRename(filePath, newName string) (bool, string) {
	err := os.Rename(filePath, newName)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

// 文件信息
func Filespec(filePath string) (bool, os.FileInfo) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, nil
	}
	return true, fileInfo
}

// 文件读取
func FileRead(filePath string) (bool, string) {
	contentByte, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return false, readErr.Error()
	}
	return true, string(contentByte)
}

// 文件分块读取(二进制)
// buffer 偏移量
// start 开始读取的位置
func FileReadBlock(filePath string, buffer int, start int) (bool, string, []byte) {
	f, err := os.Open(filePath)
	if err != nil {
		return false, err.Error(), nil
	}
	defer f.Close()

	b := make([]byte, buffer)
	n, err := f.ReadAt(b, int64(start))
	if err != nil && err != io.EOF {
		return false, err.Error(), nil
	}
	return true, "", b[:n]
}

// 文件写入
func FileWrite(filePath, content string) (bool, string) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0600)

	defer func(f io.Closer) {
		if err := f.Close(); err != nil {
			fmt.Printf("defer close file err: %v", err.Error())
		}
	}(f)

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

// 文件写入追加
func FileWriteAppend(filePath, content string) (bool, string) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)

	defer func(f io.Closer) {
		if err := f.Close(); err != nil {
			fmt.Printf("defer close file err: %v", err.Error())
		}
	}(f)

	if err != nil {
		return false, err.Error()
	} else {
		write := bufio.NewWriter(f)
		write.WriteString(content)
		write.Flush()
		return true, ""
	}
}

// 文件二进制写入
func FileWriteByte(filePath string, content []byte) (bool, string) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	defer func(f io.Closer) {
		if err := f.Close(); err != nil {
			fmt.Printf("defer close file err: %v", err.Error())
		}
	}(f)
	if err != nil {
		return false, err.Error()
	}

	var bytesBuffer bytes.Buffer
	binary.Write(&bytesBuffer, binary.LittleEndian, content)
	_, err = f.Write(bytesBuffer.Bytes())
	if err != nil {
		return false, err.Error()
	}

	return true, ""
}

// 新建文件夹
func DirMake(dirPath string) (bool, string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return false, err.Error()
	} else {
		return true, ""
	}
}

// 文件夹信息
func DirCheck(dirPath string) (bool, string, []fs.DirEntry) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "", files
}

// 文件删除
func DirDel(dirPath string) (bool, string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
