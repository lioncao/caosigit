package tools

/*
配置文件处理工具
配置文件的基本格式如下:
========================================
[global]
#我是注释行
aaa=bbb

[database]
aaa=bbb
ccc=ddd

[Service]
name=service1
type=type1

[Service]
name=service2
type=type2

=========================================
global, database 为域的名称, 不应重复, 不同的域内部是独立的命名空间
aaa=bbb 为一个参数配置,等号左边为参数名称, 右边为参数的字符串形式
以#开始的行为注释行
Service为特殊的域名称,允许重复,每个Service域均为独立的命名空间

*/

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	_CONFIG_SECTION_SERVICE = "Service"
)

type configMgr struct {
	MainData    map[string]map[string]string
	ServiceData []map[string]string
	curSection  string
}

type ConfigMgr struct {
	configMgr
}

func (this *configMgr) Load(path string) error {
	if this.MainData == nil {
		this.MainData = make(map[string]map[string]string)
	}

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	for {
		s, err1 := r.ReadString('\n')
		if err1 != nil {
			this.PraseString(s)
			break
		}
		this.PraseString(s)
	}
	return nil
}

func (this *configMgr) GetSection(s string) map[string]string {
	se, ok := this.MainData[s]
	if ok {
		return se
	} else {
		return nil
	}
}

func (this *configMgr) GetServiceData() []map[string]string {
	return this.ServiceData
}

func (this *configMgr) Get(section, key string) (string, bool) {
	s, ok := this.MainData[section]
	if ok {
		v, ok := s[key]
		return v, ok
	}
	return "", false
}

func (this *configMgr) GetInt(section, key string) (int64, bool) {
	s, ok := this.Get(section, key)
	if ok {
		i, e := strconv.ParseInt(s, 0, 64)
		if e != nil {
			return 0, false
		}
		return i, true

	}
	return 0, false
}

func (this *configMgr) GetBool(section, key string) (bool, bool) {
	s, ok := this.Get(section, key)
	if ok {
		b, e := strconv.ParseBool(s)
		if e != nil {
			return false, false
		}
		return b, true
	}
	return false, false
}

func (this *configMgr) PraseString(content string) {
	realStr := strings.TrimSpace(content)
	// realStr := strings.Trim(content, cutset)
	if realStr == "" || realStr[0] == '#' {
		return
	} else {
		if strings.ContainsAny(realStr, "[&]") {
			section := realStr[1 : len(realStr)-1]
			if strings.EqualFold(section, _CONFIG_SECTION_SERVICE) {
				sd := make(map[string]string)
				this.ServiceData = append(this.ServiceData, sd)
			} else {
				_, ok := this.MainData[section]
				if !ok {
					this.MainData[section] = make(map[string]string)
				}
			}
			this.curSection = section
		} else {
			idx := strings.Index(realStr, "=")
			if idx > 0 {
				kv := strings.SplitN(realStr, "=", 2)
				if len(kv) == 2 {
					for i := 0; i < 2; i++ {
						kv[i] = strings.TrimSpace(kv[i])
					}
					if strings.EqualFold(this.curSection, _CONFIG_SECTION_SERVICE) {
						sm := this.ServiceData[len(this.ServiceData)-1]
						//fmt.Printf("serveric push s:%s key:%s v:%s\n", this.curSection, kv[0], kv[1])
						sm[kv[0]] = kv[1]
					} else {
						datamap := this.MainData[this.curSection]
						//fmt.Printf("push s:%s key:%s v:%s\n", this.curSection, kv[0], kv[1])
						datamap[kv[0]] = kv[1]
					}
				}

			}
		}

	}
}

func (this *configMgr) EnsureInt(section, key string, defaultValue int64) int64 {
	v, ok := this.GetInt(section, key)
	if ok {
		return v
	} else {
		return defaultValue
	}
}

func (this *configMgr) EnsureBool(section, key string, defaultValue bool) bool {
	v, ok := this.GetBool(section, key)
	if ok {
		return v
	} else {
		return defaultValue
	}
}

func (this *configMgr) EnsureString(section, key string, defaultValue string) string {
	v, ok := this.Get(section, key)
	if ok {
		return v
	} else {
		return defaultValue
	}
}

func (this *configMgr) Services() []map[string]string {
	return this.ServiceData
}

// 打印所有的配置项
func (this *configMgr) Print() {
	for gname, gdata := range this.MainData {
		fmt.Println(Color(CL_YELLOW, gname))
		for key, value := range gdata {
			fmt.Println("    ", Color(CL_GREEN, key), value)
		}
		fmt.Println("")
	}

	// service部分
	if this.ServiceData != nil && len(this.ServiceData) > 0 {
		fmt.Println(Color(CL_YELLOW, "===================================\n"))
		for i, gdata := range this.ServiceData {
			str := fmt.Sprintf(Color(CL_YELLOW, "%s %d"), _CONFIG_SECTION_SERVICE, i)
			fmt.Println(str)
			for key, value := range gdata {
				fmt.Println("    ", Color(CL_GREEN, key), value)
			}
			fmt.Println("")
		}
	}

}

var config_intanse *configMgr

func GetConfigMgr() *configMgr {
	if config_intanse == nil {
		config_intanse = new(configMgr)
	}
	return config_intanse
}

func NewConfigMgr() *ConfigMgr {
	return new(ConfigMgr)
}
