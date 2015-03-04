package tools

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type weError struct {
	errorNo int
	e       string
}

func (this *weError) Error() string {
	return this.e
}

func NewError(errorNo int) error {
	return &weError{errorNo: errorNo}
}

func GetErrorNo(e error) int {
	if serr, ok := e.(*weError); ok {
		return serr.errorNo
	} else {
		return 0
	}
}

func Error(formatstring string, a ...interface{}) error {
	errstr := fmt.Sprintf(formatstring, a...)
	return errors.New(errstr)
}

//随机数
func GetRand(Max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(Max)
}

func GetTimeDay() int32 {
	t := time.Now()
	s := t.Format("20060102") //20141105格式
	port, _ := strconv.Atoi(s)
	return int32(port)
}
