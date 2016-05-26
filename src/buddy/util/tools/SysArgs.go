package tools

import (
	"strconv"
	"strings"
)

// 用于存放一组参数的数据对象
type sysArgsData struct {
	name  string
	paras []string
}

// 参数解析工具以及结果存放
type SysArgs struct {
	args  []string                // 原始的参数列表
	datas map[string]*sysArgsData // 解析后的数据存储
	usage string                  // 使用说明
	cmd   string                  // 原始命令(也就是参数列表的第0项)
}

///////////////////////////////////////////////////////////////////////////
// 设置和获取使用说明
func (this *SysArgs) SetUsage(usage string) {
	this.usage = usage
}

func (this *SysArgs) Usage() string {
	return this.usage
}

func (this *SysArgs) Cmd() string {
	return this.cmd
}

// 将传入的参数表进行解析, 并存放下来
func (this *SysArgs) Parse(args_src []string) {
	count := len(args_src)
	this.args = make([]string, count)

	if count > 0 {
		copy(this.args, args_src)
		this.cmd = this.args[0]
	}

	args := this.args

	datas := make(map[string]*sysArgsData)
	this.datas = datas

	var (
		value         string
		data, dataOrg *sysArgsData
	)

	data = nil
	for i := 1; i < count; i++ {
		value = args[i]
		bytes := []byte(value)

		if bytes[0] == '-' { // 以'-'开头
			dataOrg = datas[value]
			if dataOrg != nil { // 检查是否有重复的flag定义
				ShowWarnning("parse arg flag redefined:", i, "\""+value+"\"")
				data = dataOrg
			} else {
				data = new(sysArgsData)
				data.name = value
				data.paras = make([]string, 0)
				datas[value] = data
			}
		} else {
			if data == nil { // 前面没有参数flag
				ShowWarnning("parse arg no flag name:", i, "\""+value+"\"")
			} else {
				data.paras = append(data.paras, value)
			}
		}
	}
}

// 获取指定参数的所有后续数据
func (this *SysArgs) Values(name string) []string {
	data := this.datas[name]
	if data != nil {
		return data.paras
	}
	return nil
}

func (this *SysArgs) HasValue(name string) bool {
	data := this.datas[name]
	if data != nil {
		return true
	}
	return false
}

// 获取指定参数的单个后续数据(字符串)
func (this *SysArgs) String(name string, defaultValue string) (string, error) {
	paras := this.Values(name)
	if paras == nil || len(paras) <= 0 {
		e := Error("sys args String no para: %s", name)
		// ShowWarnning(e.Error())
		return defaultValue, e
	}
	return paras[0], nil
}

// 获取指定参数的单个后续数据(整数)
func (this *SysArgs) Int64(name string, defaultValue int64) (int64, error) {
	paras := this.Values(name)
	if paras == nil || len(paras) <= 0 {
		e := Error("sys args Int64 no para: %s", name)
		// ShowWarnning(e.Error())
		return defaultValue, e
	}
	ret, err := strconv.ParseInt(paras[0], 0, 64)
	if err != nil {
		e := Error("sys args Int64 use defaultValue: name=%s, p[0]=%s, %s", name, paras[0], err.Error())
		// ShowWarnning(e.Error())
		return defaultValue, e
	}
	return ret, nil
}

// 打印所有参数
func (this *SysArgs) Print() {

	strs := make([]string, 0)
	strs = append(strs, this.cmd)

	for name, data := range this.datas {
		msg := name
		for _, v := range data.paras {
			msg = msg + "\t" + v
		}
		strs = append(strs, msg)
	}

	ShowDebug("\n" + strings.Join(strs, "\n") + "\n")
}
