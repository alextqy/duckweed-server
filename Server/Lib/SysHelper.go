package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

func LocalIP() (bool, string, []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {

		return false, err.Error(), nil
	} else {
		var ips []string
		for _, ads := range addrs {
			if ipnet, ok := ads.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, string(ipnet.IP.String()))
				}
			}
		}
		return true, "", ips
	}
}

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

func StringToInt64(data string) (bool, string, int64) {
	res, err := strconv.ParseInt(data, 10, 64)
	if err != nil {

		return false, err.Error(), 0
	} else {
		return true, "", res
	}
}

func IntToString(data int) string {
	return strconv.Itoa(data)
}

func Int64ToString(data int64) string {
	return strconv.FormatInt(data, 10)
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

func TimeNow() time.Time {
	return time.Now()
}

func TimeNowStr() string {
	return time.Now().Format("2006-01-02 15:04:05") // 2006-01-02 15:04:05 golang立项时间
}

func TimeStamp() int64 {
	return time.Now().Unix()
}

func TimeStampMS() int64 {
	return time.Now().UnixNano()
}

func TimeStampToStr(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

func MD5(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

func EnBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func DeBase64(s64 string) (bool, string, string) {
	decoded, err := base64.StdEncoding.DecodeString(s64)
	if err != nil {

		return false, err.Error(), ""
	}
	return true, "", string(decoded)
}

func StringContains(data string, subs string) bool {
	return strings.Contains(data, subs)
}

func LogDir() string {
	return "../Log/" + strings.Split(TimeNowStr(), " ")[0] + "/"
}

func WriteLog(fileName string, content string) (bool, string) {
	if !FileExist(LogDir()) {
		b, s := DirMake(LogDir())
		if !b {
			return false, s
		}
	}
	logFile := LogDir() + fileName + ".log"
	if !FileExist(logFile) {
		b, s := FileMake(logFile)
		if !b {
			return false, s
		}
	}
	b, s := FileWriteAppend(logFile, TimeNowStr()+" "+content+""+"\n")
	if !b {
		return false, s
	}
	return true, ""
}

/*
CBC 加密
data 待加密的明文
key 秘钥
vi 向量
*/
func AesEncrypterCBC(data_s string, key_s string, iv_s string) (bool, string, string) {
	data := []byte(data_s)
	key := []byte(key_s)
	iv := []byte(iv_s)
	block, err := aes.NewCipher(key)
	if err != nil {
		return false, err.Error(), ""
	}
	padding := block.BlockSize() - len(data)%block.BlockSize()
	var paddingText []byte
	if padding == 0 {
		paddingText = bytes.Repeat([]byte{byte(block.BlockSize())}, block.BlockSize())
	} else {
		paddingText = bytes.Repeat([]byte{byte(padding)}, padding)
	}
	paddText := append(data, paddingText...)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	result := make([]byte, len(paddText))
	blockMode.CryptBlocks(result, paddText)
	return true, "", string(result)
}

/*
CBC 解密
data 待解密的密文
key 秘钥
vi 向量
*/
func AesDecrypterCBC(data_s string, key_s string, iv_s string) (bool, string, string) {
	data := []byte(data_s)
	key := []byte(key_s)
	iv := []byte(iv_s)
	block, err := aes.NewCipher(key)
	if err != nil {
		return false, err.Error(), ""
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(data))
	blockMode.CryptBlocks(result, data)
	unPadding := int(result[len(result)-1])
	return true, "", string(result[:(len(result) - unPadding)])
}

// 大小写英文字母
func RegEn(s string) bool {
	r, err := regexp.Compile("^[a-zA-Z]+$")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return r.MatchString(s)
}

// 数字
func RegNum(s string) bool {
	r, err := regexp.Compile("^[0-9]*$")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return r.MatchString(s)
}

// 中文
func RegZh(s string) bool {
	r, err := regexp.Compile("[\u4e00-\u9fa5]")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return r.MatchString(s)
}

// 英文 数字
func RegEnNum(s string) bool {
	r, err := regexp.Compile("^[a-zA-Z0-9]+$")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return r.MatchString(s)
}

// 中英文 数字 下划线 短横线 中英文(逗号 句号 分号 感叹号)
func RegWriting(s string) bool {
	r, err := regexp.Compile("^[\u4e00-\u9fa5_a-zA-Z0-9-,.;!，。；！]+$")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return r.MatchString(s)
}

// 中英文 数字 下划线 短横线
func RegAll(s string) bool {
	r, err := regexp.Compile("^[\u4e00-\u9fa5_a-zA-Z0-9-]+$")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return r.MatchString(s)
}
