package service

import (
	"buddy/util/tools"
	"encoding/xml"
	"fmt"
	"runtime"
	"strconv"
)

// 整个服务器的配置信息
type ServerConfig struct {
	ServerName           string
	Version              string
	DebugMode            bool
	IP                   string
	Port                 string
	Services             []ServiceInfo `xml:"Service"`
	GOMAXPROCS           int           // 允许使用的最大cpu数量
	CommonConfigFileName string        // 通用配置所在文件
}

// 每个service的配置信息
type ServiceInfo struct {
	XMLName     xml.Name `xml:"Service"`
	Type        string
	Name        string
	ConfigFile  string
	Params      string
	DebugMode   bool
	Description string
	Status      int32
}

// 获取指定Service的配置参数
func (this *ServerConfig) GetServiceInfo(serviceName string) *ServiceInfo {

	for _, v := range this.Services {
		if v.Name == serviceName {
			return &v
		}
	}
	return nil
}

// 获取指定Service的配置参数
func (this *ServerConfig) GetServiceStatus(serviceName string) int32 {

	for _, v := range this.Services {
		if v.Name == serviceName {
			return v.Status
		}
	}
	return CLOSE
}

// 打印配置信息
func (this *ServerConfig) PrintConfigInfo() {
	var st string
	var debugDesc string

	// 服务器名称以及版本号
	tools.ShowInfo(""+this.ServerName, "(", this.Version, ")")

	// 服务器监听情况
	tools.ShowInfo("Listen Addr\t", this.IP+":"+this.Port)

	// 调试模式开关
	if this.DebugMode {
		st = "Open"
	} else {
		st = "Close"
	}
	tools.ShowInfo("DebugMode\t", st)

	// CPU 相关信息
	tools.ShowInfo("LOC_NUMCPU\t", runtime.NumCPU())
	tools.ShowInfo("GOMAXPROCS\t", this.GOMAXPROCS)

	fuckPrint := false

	fmtStr := ("Service%6s %20s %15s %80s %10s %s")
	// 各项服务器信息
	if fuckPrint {
		tools.ShowInfo(fmt.Sprintf(fmtStr, "Index", "Key", "Status", "ConfigFile", "Params", "Description"))
	} else {
		tools.ShowInfo("Service\t", "Index\t", "Key\t\t\t", "Status\t", "ConfigFile\t", "Params\t", "Description\t")
	}
	var count int32 = 0
	for k, v := range this.Services {

		if v.Status == 1 {
			st = tools.Color(tools.CL_GREEN, "Open")
			count++
		} else {
			st = tools.Color(tools.CL_GRAY, "Close")
		}

		if v.DebugMode {
			debugDesc = tools.Color(tools.CL_RED, "(debug)")
		} else {
			debugDesc = tools.Color(tools.CL_GREEN, "(release)")
		}
		if fuckPrint {

			tools.ShowInfo(fmt.Sprintf(fmtStr, strconv.FormatInt(int64(k), 10), v.Name+debugDesc, st, v.ConfigFile, v.Params, v.Description))
		} else {
			tools.ShowInfo("Service\t", k, "\t", v.Name+debugDesc, "\t", st+"\t", v.ConfigFile+"\t", v.Params+"\t", v.Description)
		}
	}
	tools.ShowInfo("Service Count:\t", count, "\t\t\t")

}

// 装载服务器配置
func LoadServerConfig(fileName string) (*ServerConfig, error) {
	if fileName == "" {
		tools.ShowInfo("Config file name is nil")
		return nil, tools.Error("config file name is nil")
	}
	cf := new(ServerConfig)
	err := tools.LoadXmlFile(fileName, cf)

	if err != nil {
		return nil, err
	} else {
		return cf, nil
	}
}
