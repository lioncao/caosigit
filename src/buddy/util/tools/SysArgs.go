package tools

import (
	"strconv"
)

type sysArgsData struct {
	name  string
	paras []string
}

type SysArgs struct {
	datas map[string]*sysArgsData
	usage string
}

func (this *SysArgs) SetUsage(usage string) {
	this.usage = this.usage
}

func (this *SysArgs) Usage() string {
	return this.usage
}

func (this *SysArgs) Parse(args []string) {
	datas := make(map[string]*sysArgsData)
	this.datas = datas

	var (
		value         string
		data, dataOrg *sysArgsData
	)

	count := len(args)
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

func (this *SysArgs) Values(name string) []string {
	data := this.datas[name]
	if data != nil {
		return data.paras
	}
	return nil
}

func (this *SysArgs) String(name string, defaultValue string) (string, error) {
	paras := this.Values(name)
	if paras == nil || len(paras) <= 0 {
		e := Error("sys args String no para: %s", name)
		// ShowWarnning(e.Error())
		return defaultValue, e
	}
	return paras[0], nil
}

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

func (this *SysArgs) Print() {

	for name, data := range this.datas {
		msg := name
		for _, v := range data.paras {
			msg = msg + "\t" + v
		}
		ShowDebug(msg)
	}
}
