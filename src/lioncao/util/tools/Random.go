package tools

////////////////////////////////////////////////////////////////////////////////////////////////////
// tools包为基础工具包, 不应依赖任何自建的package
////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

// md5生成工具
// 结果为32小写md5字符串
func Md5(str string) (string, error) {
	h := md5.New()

	_, err := io.WriteString(h, str)
	if err != nil {
		return "", err
	}

	strMd5 := fmt.Sprintf("%x", h.Sum(nil))

	return strMd5, nil
}

// 在[a,b]区间获取随机整数
func Random(a, b int64) (random int64, err error) {
	if a == b {
		return a, nil
	}
	if a > b {
		a, b = b, a
	}
	delta := b - a + 1
	max := big.NewInt(delta)

	// rand.Read([]byte(TimeString(0)))
	i, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	rnd := a + i.Int64()
	return rnd, nil
}

// 获取随机字符串
func RandomStringBase(size int, srcChars []byte, srcSize int) string {
	if size <= 0 || srcChars == nil {
		return ""
	}

	if srcSize <= 0 {
		srcSize = len(srcChars)
	}

	buf := make([]byte, size+1)
	var (
		max int64 = int64(srcSize - 1)
		rnd int64
	)

	for i := 0; i < size; i++ {
		rnd, _ = Random(0, max)
		buf[i] = srcChars[rnd]
	}
	// buf[size] = 0

	return string(buf[:size])
}

var hexStringChars []byte = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f'}
var hexStringCharsSize int = len(hexStringChars)

// 获取随机16进制字符串(0-9, a-f)
func RandomStringHex(size int) string {
	return RandomStringBase(size, hexStringChars, hexStringCharsSize)
}

var rndStringChars []byte = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
var rndStringCharsSize int = len(rndStringChars)

// 获取随机字符串(a-z , 0-9)
func RandomString(size int) string {
	return RandomStringBase(size, rndStringChars, rndStringCharsSize)
}

const (
	INDEX_TIME_FMT = "060102150405" // 创建index时的时间格式
)

func CreatRandomIndex(pre string, privateNum int64, userid int64) string {
	timeStr := TimeStringFmt(EMPTY_TIME, INDEX_TIME_FMT)
	return fmt.Sprintf("%s_%d_%s_%d", pre, userid, timeStr, privateNum)
}
