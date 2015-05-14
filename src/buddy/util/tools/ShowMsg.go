package tools

import (
	"fmt"
	"strconv"
	"time"
)

const (
	UTF8_BOM     = "\xEF\xBB\xBF"
	UTF8_BOM_LEN = len(UTF8_BOM)
)

func CheckUTF8_BOM(data []byte) bool {
	if data == nil {
		return false
	}
	length := len(data)
	if length < UTF8_BOM_LEN {
		return false
	}

	if string(data[0:UTF8_BOM_LEN]) == UTF8_BOM {
		return true
	}
	return false
}

// 输出文字颜色控制相关常数
const (
	CL_RESET = "\033[0m"
	CL_CLS   = "\033[2J"
	CL_CLL   = "\033[K"

	// font settings
	CL_BOLD   = "\033[1m"
	CL_NORM   = CL_RESET
	CL_NORMAL = CL_RESET
	CL_NONE   = CL_RESET
	// foreground color and bold font (bright color on windows)
	CL_WHITE   = "\033[1;37m"
	CL_GRAY    = "\033[1;30m"
	CL_RED     = "\033[1;31m"
	CL_GREEN   = "\033[1;32m"
	CL_YELLOW  = "\033[1;33m"
	CL_BLUE    = "\033[1;34m"
	CL_MAGENTA = "\033[1;35m"
	CL_CYAN    = "\033[1;36m"

	// background color
	CL_BG_BLACK   = "\033[40m"
	CL_BG_RED     = "\033[41m"
	CL_BG_GREEN   = "\033[42m"
	CL_BG_YELLOW  = "\033[43m"
	CL_BG_BLUE    = "\033[44m"
	CL_BG_MAGENTA = "\033[45m"
	CL_BG_CYAN    = "\033[46m"
	CL_BG_WHITE   = "\033[47m"
	// foreground color and normal font (normal color on windows)
	CL_LT_BLACK   = "\033[0;30m"
	CL_LT_RED     = "\033[0;31m"
	CL_LT_GREEN   = "\033[0;32m"
	CL_LT_YELLOW  = "\033[0;33m"
	CL_LT_BLUE    = "\033[0;34m"
	CL_LT_MAGENTA = "\033[0;35m"
	CL_LT_CYAN    = "\033[0;36m"
	CL_LT_WHITE   = "\033[0;37m"
	// foreground color and bold font (bright color on windows)
	CL_BT_BLACK   = "\033[1;30m"
	CL_BT_RED     = "\033[1;31m"
	CL_BT_GREEN   = "\033[1;32m"
	CL_BT_YELLOW  = "\033[1;33m"
	CL_BT_BLUE    = "\033[1;34m"
	CL_BT_MAGENTA = "\033[1;35m"
	CL_BT_CYAN    = "\033[1;36m"
	CL_BT_WHITE   = "\033[1;37m"

	CL_WTBL = "\033[37;44m"   // white on blue
	CL_XXBL = "\033[0;44m"    // default on blue
	CL_PASS = "\033[0;32;42m" // green on green

	CL_SPACE = "           " // space aquivalent of the print messages
)

// 给文字添加颜色
// e.g:
//		Color(CL_YELLOW , "WARNNING")
func Color(colorStr string, srcStr string) string {
	return colorStr + srcStr + CL_RESET
}

const (
	// 打印开关
	flag_SHOW_INFO     = 0x1
	flag_SHOW_DEBUG    = 0x2
	flag_SHOW_WARNNING = 0x4
	flag_SHOW_ERROR    = 0x8
	// 时间格式
	// TIME_FMT = "\033[32m[06-01-02 15:04:05.000]\033[0m"
	TIME_FMT           = "[06-01-02 15:04:05.000]"
	TIME_FMT_DIGIT_DAY = "060102"
	TIME_FMT_DIGIT_SEC = "060102030405"
	// 信息题头
	SHOW_TITLE_INFO     = CL_GREEN + "[INFO]" + CL_RESET
	SHOW_TITLE_DEBUG    = CL_BLUE + "[DEBUG]" + CL_RESET
	SHOW_TITLE_WARNNING = CL_YELLOW + "[WARN]" + CL_RESET
	SHOW_TITLE_ERROR    = CL_RED + "[ERROR]" + CL_RESET

	EMPTY_TIME = -1
)

// 打印标记
var showFlags int64 = flag_SHOW_INFO | flag_SHOW_DEBUG | flag_SHOW_WARNNING | flag_SHOW_ERROR

func SetShowFlag(flags int64) {
	showFlags = flags
}

func TimeString(timeValue int64) string {
	return TimeStringFmt(timeValue, TIME_FMT)
}

func TimeStringFmt(timeValue int64, timeFmt string) string {
	var t time.Time
	if timeValue != EMPTY_TIME {
		t = time.Unix(timeValue, 0)
	} else {
		t = time.Now()
	}
	return t.Format(timeFmt)
}

func TimeDigitValue(sec int64, timeFmt string) int64 {
	s := TimeStringFmt(sec, timeFmt)

	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return d
}

func ShowInfo(a ...interface{}) {
	if (showFlags & flag_SHOW_INFO) != 0 {
		fmt.Println(SHOW_TITLE_INFO, TimeString(EMPTY_TIME), a)
	}
}
func ShowDebug(a ...interface{}) {
	if (showFlags & flag_SHOW_DEBUG) != 0 {
		fmt.Println(SHOW_TITLE_DEBUG, TimeString(EMPTY_TIME), a)
	}
}

func ShowWarnning(a ...interface{}) {
	if (showFlags & flag_SHOW_WARNNING) != 0 {
		fmt.Println(SHOW_TITLE_WARNNING, TimeString(EMPTY_TIME), a)
	}
}

func ShowError(a ...interface{}) {
	if (showFlags & flag_SHOW_ERROR) != 0 {
		fmt.Println(SHOW_TITLE_ERROR, TimeString(EMPTY_TIME), a)
	}
}
func CaoSiShowDebug(a ...interface{}) {
	if (showFlags & flag_SHOW_DEBUG) != 0 {
		fmt.Println(CL_CYAN+"[CAOSI_DEBUG]"+CL_RESET, TimeString(EMPTY_TIME), a)
	}
}
