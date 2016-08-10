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
	if showUseColor {
		return colorStr + srcStr + CL_RESET
	} else {
		return srcStr
	}
}

// 不受开关控制的颜色工具函数
func ColorForce(colorStr string, srcStr string) string {
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
	show_TITLE_INFO_BASE     = "[INFO]"
	show_TITLE_DEBUG_BASE    = "[DBUG]"
	show_TITLE_WARNNING_BASE = "[WARN]"
	show_TITLE_ERROR_BASE    = "[ERRO]"

	show_TITLE_INFO_WITH_COLOR     = CL_GREEN + show_TITLE_INFO_BASE + CL_RESET
	show_TITLE_DEBUG_WITH_COLOR    = CL_BLUE + show_TITLE_DEBUG_BASE + CL_RESET
	show_TITLE_WARNNING_WITH_COLOR = CL_YELLOW + show_TITLE_WARNNING_BASE + CL_RESET
	show_TITLE_ERROR_WITH_COLOR    = CL_RED + show_TITLE_ERROR_BASE + CL_RESET

	EMPTY_TIME = -1
)

// 打印标记
var (
	showFlags         int64 = flag_SHOW_INFO | flag_SHOW_DEBUG | flag_SHOW_WARNNING | flag_SHOW_ERROR
	showFlag_info           = true
	showFlag_debug          = true
	showFlag_warnning       = true
	showFlag_error          = true

	show_TITLE_INFO     = show_TITLE_INFO_WITH_COLOR
	show_TITLE_DEBUG    = show_TITLE_DEBUG_WITH_COLOR
	show_TITLE_WARNNING = show_TITLE_WARNNING_WITH_COLOR
	show_TITLE_ERROR    = show_TITLE_ERROR_WITH_COLOR

	showUseColor = true // 是否使用颜色打印
)

// 设置是否在showMsg中使用颜色
// useColor:  true 表示使用颜色， 否色不使用
func SetShowMsgColorFlag(useColor bool) {
	showUseColor = useColor
	if !useColor {
		show_TITLE_INFO = show_TITLE_INFO_BASE
		show_TITLE_DEBUG = show_TITLE_DEBUG_BASE
		show_TITLE_WARNNING = show_TITLE_WARNNING_BASE
		show_TITLE_ERROR = show_TITLE_ERROR_BASE
	} else {
		show_TITLE_INFO = show_TITLE_INFO_WITH_COLOR
		show_TITLE_DEBUG = show_TITLE_DEBUG_WITH_COLOR
		show_TITLE_WARNNING = show_TITLE_WARNNING_WITH_COLOR
		show_TITLE_ERROR = show_TITLE_ERROR_WITH_COLOR
	}
}
func GetShowMsgColorFlag() bool {
	return showUseColor
}

func SetShowFlag(flags int64) {
	showFlags = flags

	if (showFlags & flag_SHOW_INFO) != 0 {
		showFlag_info = true
	} else {
		showFlag_info = false
	}

	if (showFlags & flag_SHOW_DEBUG) != 0 {
		showFlag_debug = true
	} else {
		showFlag_debug = false
	}

	if (showFlags & flag_SHOW_WARNNING) != 0 {
		showFlag_warnning = true
	} else {
		showFlag_warnning = false
	}

	if (showFlags & flag_SHOW_ERROR) != 0 {
		showFlag_error = true
	} else {
		showFlag_error = false
	}
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
	if showFlag_info {
		fmt.Println(show_TITLE_INFO, TimeString(EMPTY_TIME), a)
	}
}
func ShowInfoF(fmtStr string, a ...interface{}) {
	if showFlag_info {
		fmt.Println(show_TITLE_INFO, TimeString(EMPTY_TIME), fmt.Sprintf(fmtStr, a...))
	}
}
func ShowDebug(a ...interface{}) {
	if showFlag_debug {
		fmt.Println(show_TITLE_DEBUG, TimeString(EMPTY_TIME), a)
	}
}

func ShowDebugF(fmtStr string, a ...interface{}) {
	if showFlag_debug {
		fmt.Println(show_TITLE_DEBUG, TimeString(EMPTY_TIME), fmt.Sprintf(fmtStr, a...))
	}
}

func ShowWarnning(a ...interface{}) {
	if showFlag_warnning {
		fmt.Println(show_TITLE_WARNNING, TimeString(EMPTY_TIME), a)
	}
}
func ShowWarnningF(fmtStr string, a ...interface{}) {
	if showFlag_warnning {
		fmt.Println(show_TITLE_WARNNING, TimeString(EMPTY_TIME), fmt.Sprintf(fmtStr, a...))
	}
}

func ShowError(a ...interface{}) {
	if showFlag_error {
		fmt.Println(show_TITLE_ERROR, TimeString(EMPTY_TIME), a)
	}
}

func ShowErrorF(fmtStr string, a ...interface{}) {
	if showFlag_error {
		fmt.Println(show_TITLE_ERROR, TimeString(EMPTY_TIME), fmt.Sprintf(fmtStr, a...))
	}
}

func CaoSiShowDebug(a ...interface{}) {
	if showFlag_debug {
		fmt.Println(Color(CL_CYAN, "[CAOSI_DEBUG]"), TimeString(EMPTY_TIME), a)
	}
}
