package tools

import (
	"math/rand"
	"strconv"
	"time"
)

//自定义错误格式
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

// 随机数
// 在 [0, Max)的左闭右开区间中取随机数
func GetRand(Max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(Max)
}

// 随机数2
// 在 [min, max]的闭区间中取随机数
func GetRand2(min int, max int) int {
	if min > max {
		min, max = max, min
	}
	return GetRand(max-min+1) + min
}

func GetTimeDay() int32 {
	t := time.Now()
	s := t.Format("20060102") //20141105格式
	port, _ := strconv.Atoi(s)
	return int32(port)
}

func GetDayByTime(t time.Time) int32 {
	s := t.Format("20060102") //20141105格式
	port, _ := strconv.Atoi(s)
	return int32(port)
}

func GetTimeMonth() int32 {
	t := time.Now()
	s := t.Format("200601") //20141105格式
	port, _ := strconv.Atoi(s)
	return int32(port)
}
