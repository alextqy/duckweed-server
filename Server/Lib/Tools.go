package Lib

import (
	"bytes"
	"encoding/binary"
	"strconv"
)

func StringToByte(data string) []byte {
	return []byte(data)
}

func ByteToString(data []byte) string {
	return string(data)
}

func StringToInt(data string) (bool, string, int) {
	res, err := strconv.Atoi(data)
	if err != nil {
		return false, err.Error(), 0
	} else {
		return true, "", res
	}
}

func IntToString(data int) string {
	return strconv.Itoa(data)
}

func StringToFloat64(data string) (bool, string, float64) {
	s, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return false, err.Error(), 0
	} else {
		return true, "", s
	}
}

func Float64ToString(data float64) string {
	return strconv.FormatFloat(data, 'E', -1, 32)
}

func IntToBytes(data int) []byte {
	dataInt := int32(data)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, dataInt)
	return bytesBuffer.Bytes()
}

func BytesToInt(data []byte) int {
	bytesBuffer := bytes.NewBuffer(data)
	var dataInt int32
	binary.Read(bytesBuffer, binary.BigEndian, &dataInt)
	return int(dataInt)
}
