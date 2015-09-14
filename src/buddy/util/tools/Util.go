package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func RuntimeInfo() (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(2)
}

// http发送简单文本
func HttpSendSimplePage(w *http.ResponseWriter, code int, content string) {
	if w != nil {
		(*w).WriteHeader(code)
		(*w).Write([]byte(content))
	}
}

// 发送空返回
func HttpSend200Empty(w *http.ResponseWriter) {
	HttpSendSimplePage(w, 200, "empty")
}

// 发送404页面
func HttpSend404NotFound(w *http.ResponseWriter) {
	HttpSendSimplePage(w, 404, "NOT FOUND")
}

func Error(fmtStr string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(fmtStr, a))
}

const (
	// 拼接session字符串时用到的时间格式
	SESSION_STR_TIME_FMT = "060102150405"
)

// 获取一个随机的sessionid
func MakeSessionId(title string) string {
	return fmt.Sprintf("%s_%s_%s", title, TimeStringFmt(EMPTY_TIME, SESSION_STR_TIME_FMT), RandomString(16))
}

// 统一的redisKey生成函数
// dbInde：			数据库id（实际功能为给数据划分区段）
// serviceKey：		服务编号,对应唯一的一个service, 例如: auth_001, guobao_003 , balala_100
// userId：			用户的唯一id
// srcKey:			原始key（即功能自己用来区分的key）
func MakeRedisKey(dbIndex int64, serviceKey string, userId int64, srcKey string) string {
	return fmt.Sprintf("%d/%s/%d/%s", dbIndex, serviceKey, userId, srcKey)
}

// 将lua中的table的字符串形式转化为go中的数据结构
func DecodeLuaTableString(name string, tableStr string, data interface{}) error {
	str := strings.TrimSpace(tableStr)
	if str == "" || str == "{}" {
		return nil
	}

	str = strings.Replace(str, "{", "[", -1)
	str = strings.Replace(str, "}", "]", -1)
	str = fmt.Sprintf("{\"%s\":%s}", name, str)
	err := json.Unmarshal([]byte(str), data)
	if err != nil {
		ShowError(name, err.Error(), tableStr, str)
	}
	return err
}
