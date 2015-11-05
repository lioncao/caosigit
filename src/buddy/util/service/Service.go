package service

import (
	bdsocket "buddy/net/socket"
	"buddy/util/cmd"
	"buddy/util/tools"
	"fmt"
	"net"
	"os"
	"time"
)

// 服务器的开启关闭配置常数
const (
	OPEN               = 1
	CLOSE              = 0
	DEFAULT_HEART_BEAT = time.Millisecond * 50

	CFG_DOMAIN_GLOBAL = "global"
)

// 服务器的运行状态
const (
	RUNTIME_STATE_DEFAULT = iota // 默认状态(尚未初始化完毕)
	RUNTIME_STATE_RUNNING        // 运行中
	RUNTIME_STATE_DESTROY        // 已经关闭
)

type TcpCallbacks struct {
	Session interface{}
}

// tcp连接 MessageHandler的三个回调, 也就是bdsocket.TcpCallbacks interface{}的实现
func (this *TcpCallbacks) SetSessionData(data interface{}) {
	this.Session = data
}
func (this *TcpCallbacks) GetSessionData() interface{} {
	return this.Session
}

func (this *TcpCallbacks) TcpOnReceiveCallback(handler *bdsocket.MessageHandler, buf []byte, bufSize int) error {
	tools.CaoSiShowDebug("service tcp call backs TcpOnReceiveCallback")
	return nil
}
func (this *TcpCallbacks) TcpDoMessageCallback(handler *bdsocket.MessageHandler) error {
	tools.CaoSiShowDebug("service tcp call backs TcpDoMessageCallback")
	return nil
}
func (this *TcpCallbacks) TcpOnWriteCallback(handler *bdsocket.MessageHandler) error {
	tools.CaoSiShowDebug("service tcp call backs TcpOnWriteCallback")
	return nil
}

// service 接口定义,每个具体的service均应该实现此接口
type Service interface {
	SetInfo(cf ServiceInfo) error                              // 设置信息
	GetInfo() *ServiceInfo                                     // 获取信息
	Run() error                                                // 启动
	Stop() error                                               // 结束
	Reset() error                                              // 复位(停止服务,并回复到初始状态)
	DoCmd(session *cmd.CmdSession, cmdData *cmd.CmdData) error // 命令行响应
	SetRuntimeState(state int) int                             // 状态设置
	GetRuntimeState() int                                      // 状态检查
	OnSig(sig os.Signal)                                       // 系统信号处理
}

// service实现的基类
type ServiceSuper struct {
	ServiceInfo   ServiceInfo
	ServiceConfig tools.CommonJsonData
	TitleOnShow   string
	RuntimeState  int
	HeartBeat     time.Duration
	DebugMod      bool   // 调试开关
	ServerId      int64  // 服务器唯一id
	ServerKey     string // 服务器key __Q:可能不需要这个东西
	Language      string // 语言代码

	TcpCB bdsocket.TcpCallbacks
}

////////////////////////////////////////////////////////////////////////////////
// Service interface 的实现 begin
////////////////////////////////////////////////////////////////////////////////
// 设置信息
func (this *ServiceSuper) SetInfo(cf ServiceInfo) error {
	this.ServiceInfo = cf
	this.TitleOnShow = cf.Name
	this.DebugMod = cf.DebugMode
	if this.DebugMod {
		tools.ShowWarnning("=================================================")
		tools.ShowWarnning("=================================================")
		tools.ShowWarnning(" This is ", cf.Name)
		tools.ShowWarnning(" DebugMod is running")
		tools.ShowWarnning("=================================================")
		tools.ShowWarnning("=================================================")
	}

	return nil
}

// 获取信息
func (this *ServiceSuper) GetInfo() *ServiceInfo {
	return &this.ServiceInfo
}

// 启动
func (this *ServiceSuper) Run() error {

	return nil
}

// 结束
func (this *ServiceSuper) Stop() error {
	return nil
}

// 复位
func (this *ServiceSuper) Reset() error {
	return nil
}

// 命令行响应
func (this *ServiceSuper) DoCmd(session *cmd.CmdSession, cmdData *cmd.CmdData) error {
	return nil
}

// 状态设置
func (this *ServiceSuper) SetRuntimeState(state int) int {
	state, this.RuntimeState = this.RuntimeState, state
	return state
}

// 状态检查
func (this *ServiceSuper) GetRuntimeState() int {
	return this.RuntimeState
}

////////////////////////////////////////////////////////////////////////////////
// Service interface 的实现 end
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////
// ServiceSuper 基本功能实现 start
////////////////////////////////////////////////////////////////////////////////////////////////////

func (this *ServiceSuper) SetTcpCB(tcpCB bdsocket.TcpCallbacks) {
	this.TcpCB = tcpCB
}

// 启动socket监听服务
func (this *ServiceSuper) StartSocketService(ip, port string) error {
	go bdsocket.StartTcp(ip, port, this.DoTcp)
	return nil
}

// tcp逻辑处理
func (this *ServiceSuper) DoTcp(conn net.Conn) error {

	var err error
	//var length int
	var handler *bdsocket.MessageHandler

	// 创建消息处理句柄并设置参数
	handler, err = bdsocket.CreateMessageHandler(conn)
	if err != nil {
		tools.ShowError("CreateMessageHandler err,", conn, err)
		return nil
	}

	// 设置回调
	handler.SetCallback(this.TcpCB)

	// 启动
	err = handler.RunTcp(this.HeartBeat)
	if err != nil {
		pc, file, line, ok := tools.RuntimeInfo()
		tools.ShowError(pc, file, line, ok, "do tcp err :", err.Error())
	}
	return nil
}

func (this *ServiceSuper) OnSig(sig os.Signal) {
	info := this.GetInfo()
	tools.ShowInfo("ServiceSuper: service", info.Name, "OnSig", sig)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// ServiceSuper 基本功能实现 end
////////////////////////////////////////////////////////////////////////////////////////////////////

type ServiceCreateFunc func(sType string) (*Service, error)

// 服务器中的service管理器
type ServiceManager struct {
	Conf       *ServerConfig       // 全服的当前配置信息
	Factory    ServiceCreateFunc   // Service的工厂方法
	sMap       map[string]*Service // 实际service列表
	cmdMachine *cmd.CmdMachine
}

func (this *ServiceManager) Init(factoryFunc ServiceCreateFunc) error {
	if factoryFunc == nil {
		return tools.Error("factoryFunc is nil,init ServiceManager faild")
	}

	this.sMap = make(map[string]*Service)
	this.Conf = new(ServerConfig)
	this.Factory = factoryFunc
	return nil
}

// func (this *ServiceManager) LoadConfig(cfgName string) error {
// 	var err error
// 	this.Conf, err = LoadServerConfig(cfgName)
// 	return err
// }

// 从通用配置数据中初始化配置信息
func (this *ServiceManager) InitConfig(cfgMgr *tools.ConfigMgr) error {
	var (
		info *ServiceInfo
	)

	this.Conf.ServerName = cfgMgr.EnsureString(CFG_DOMAIN_GLOBAL, "server_name", "")
	this.Conf.DebugMode = cfgMgr.EnsureBool(CFG_DOMAIN_GLOBAL, "debug_mode", false)
	this.Conf.GOMAXPROCS = int(cfgMgr.EnsureInt(CFG_DOMAIN_GLOBAL, "cpunum", 0))
	services := cfgMgr.GetServiceData()
	for _, service_params := range services {
		info = new(ServiceInfo)
		info.Init(service_params)
		this.Conf.Services = append(this.Conf.Services, info)
	}

	return nil
}

// 启动各项服务
func (this *ServiceManager) StartServices(factoryFunc ServiceCreateFunc) error {

	// 确定工厂方法
	facF := factoryFunc
	if facF == nil {
		facF = this.Factory
	}

	if facF == nil {
		return tools.Error("Factory is nil")
	}

	// 准备遍历Service列表
	list := this.Conf.Services
	sNum := len(list)

	var s *Service
	var sInfo *ServiceInfo
	var err error
	var k string

	// 遍历服务器信息
	for i := 0; i < sNum; i++ {
		sInfo = list[i]

		k = sInfo.Name
		s = this.Get(k)

		// 从表中查找当前的信息
		if s != nil {
			tools.ShowWarnning(fmt.Sprintf("service %s was exist when start", k))
			(*s).Reset()     // 对已经存在的service进行重置
			this.Set(k, nil) // 清理原有的对象
			s = nil
		} else {
			if sInfo.Status == OPEN {
				s, err = facF(sInfo.Type)
				if s == nil || err != nil {
					tools.ShowError("service factory err:", k, err)
				}
			}
		}

		if s != nil && sInfo.Status == OPEN {
			// 启动有效Service 的goroutin
			(*s).SetInfo(*sInfo)
			this.Set(k, s)
			tools.ShowInfo(fmt.Sprintf("service run type=%s , name = %s", sInfo.Type, sInfo.Name))
			go (*s).Run()
		}
	}

	return nil
}

// 获取指定的Service句柄
func (this *ServiceManager) Get(key string) *Service {
	if key == "" {
		return nil
	}

	s, ok := this.sMap[key]
	if !ok || s == nil {
		return nil
	}
	return s
}

// 添加一个Service到ServerManager中
// s == nil 表示删除
func (this *ServiceManager) Set(key string, s *Service) (*Service, error) {
	if key == "" {
		return nil, tools.Error("ServiceInfo.Set() key is empty")
	}

	ts, _ := this.sMap[key]
	this.sMap[key] = s
	return ts, nil
}

func (this *ServiceManager) StartCmd(ip, port, password, title string) error {
	cm := cmd.NewCmdMachine(password, title, this.DoCmd)
	this.cmdMachine = cm
	go bdsocket.StartTcp(ip, port, cm.DoTcp)
	return nil
}

func (this *ServiceManager) DoCmd(session *cmd.CmdSession, cmdData *cmd.CmdData) error {
	return nil
}

func (this *ServiceManager) OnSig(sig os.Signal) {
	for _, service := range this.sMap {
		(*service).OnSig(sig)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// service status
// 每个service的运行状态信息
////////////////////////////////////////////////////////////////////////////////////////////////////
type ServiceStatus struct {
	Id     int64  // 服务统一编号
	Type   string // 服务类型
	Name   string // 服务名称
	Status int64  // 服务当前运行状态
	Ip     string // 连接ip
	Port   int32  // 连接端口
	// Sessions map[int]*util.SimpleSessionInfo // 所有在线session简要信息
}
