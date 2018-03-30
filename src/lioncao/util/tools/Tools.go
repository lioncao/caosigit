package tools

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

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

func LoadXmlFile(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	return nil
}

func LoadJsonFile(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	return nil
}
