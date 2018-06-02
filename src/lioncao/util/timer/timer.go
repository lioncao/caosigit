package timer

import (
	ga "3rdparty/goArrayList/goArrayList"
	"lioncao/util/tools"
	"time"
)

/**
TimerManager 使用的注意事项
    1. 一个TimerManager对象只服务于一个线程,不允许多个线程使用同一个,
      其具体表现就是, 只能在一个线程中调用这个对象的AddTimer函数,
      并应在其主循环中调用OnMainLoop()函数
    2. para参数请尽量使用数值或者指针, 尽量避免对象的传递

*/

const (
	TEMP_CACHE_INIT_CAPACITY   = 10
	TIMER_LIST_INIT_CAPACITY   = 100
	DO_FUNC_LIST_INIT_CAPACITY = 20
)

func _unix_time_ms(t *time.Time) int64 {
	return (*t).UnixNano() / 1000000
}
func _unix_time_ms_now() int64 {
	now := time.Now()
	return _unix_time_ms(&now)
}

// t ~ 执行timer时的unix毫秒值
// count ~ 当前是第几次执行
// para ~ 被传递的参数
type TimerFunc func(t int64, count int64, para interface{})

type timerHandler struct {
	Id int64 // 句柄id

	CreatTime  int64 // 创建时间
	DelayTime  int64 // 首次延迟(ms)
	DeltaTime  int64 // 执行间隔(ms)
	NextDoTime int64 // 下次执行时间

	Repeat int64 // 重复次数
	Count  int64 // 已执行次数

	Func TimerFunc // 方法
	Para interface{}

	Deleted bool

	mgr *TimerManager
}

type timerDo struct {
	DoTime int64     // 执行时间
	Count  int64     // 当前是第几次执行
	Func   TimerFunc // 方法
	Para   interface{}
}

type TimerManager struct {
	TempLock  tools.FastLock
	TempCache *ga.ArrayList // 用于给外部添加的临时缓存

	ListLock  tools.FastLock
	TimerList *ga.ArrayList // 当前正在生效的timer列表
	TimerMap  map[int64]*timerHandler

	// 需要执行的列表
	DoFuncLock tools.FastLock
	DoFuncList []*timerDo

	NextHandlerId int64
}

/******************************************************************************
	TimerManager funcs
******************************************************************************/
func NewTimerManager() *TimerManager {
	this := new(TimerManager)
	this.TempCache = ga.ArrayListNew(TEMP_CACHE_INIT_CAPACITY)
	this.TimerList = ga.ArrayListNew(TIMER_LIST_INIT_CAPACITY)
	this.TimerMap = make(map[int64]*timerHandler)

	this.DoFuncList = make([]*timerDo, 0, DO_FUNC_LIST_INIT_CAPACITY)
	return this
}

// 启动timer
func (this *TimerManager) Start() {
	go this.Run()
}

func (this *TimerManager) Run() {
	for {
		t := _unix_time_ms(nil)
		this.onUpdate(t)
	}
}

// t: 当前时间的毫秒值
func (this *TimerManager) onUpdate(t int64) {
	var (
		tmpList *ga.ArrayList
		handler *timerHandler
		delList []int
	)

	///////////////////////////////////////////////////////////////////////////
	// 先处理当前的列表
	delList = make([]int, 0, 20) // 用于缓存需要删除的下标
	cnt := this.TimerList.Size()
	for i := 0; i < cnt; i++ {
		handler = this.TimerList.Get(i).(*timerHandler)
		handler.onUpdate(t)
		if !handler.valid() {
			delList = append(delList, i)      // 列表中缓存需要删除的下标
			delete(this.TimerMap, handler.Id) // map中直接删除
		}
	}

	if len(delList) > 0 {
		// 倒序删除
		size := len(delList)
		for i := size - 1; i >= 0; i-- {
			this.TimerList.Remove(delList[i])
		}
	}

	///////////////////////////////////////////////////////////////////////////
	// 接下来处理临时列表中的数据
	// 取出临时表
	tmpList = nil
	if this.TempCache.Size() > 0 {
		for this.TempLock.Lock() {
			defer this.TempLock.Unlock()
			tmpList = this.TempCache
			this.TempCache = ga.ArrayListNew(TEMP_CACHE_INIT_CAPACITY)
		}
	}

	if tmpList != nil {
		size := tmpList.Size()
		for i := 0; i < size; i++ {
			handler = (tmpList.Get(i)).(*timerHandler)
			handler.onUpdate(t) // 在这里尝试执行一次
			if handler.valid() {
				// 插入到主列表中
				this.handlerInsertToTimerList(handler)
			} else {
				this.RemoveTimer(handler.Id)
			}

		}
	}

}

func (this *TimerManager) handlerInsertToTimerList(newHandler *timerHandler) {
	var (
		handler *timerHandler
	)
	size := this.TimerList.Size()
	for i := 0; i < size; i++ {
		handler = (this.TimerList.Get(i)).(*timerHandler)
		if newHandler.NextDoTime < handler.NextDoTime {
			this.TimerList.Insert(i, newHandler)
			return
		}
	}

	this.TimerList.Append(newHandler)

}

func (this *TimerManager) PopDoFuncs() []*timerDo {

	for this.DoFuncLock.Lock() {
		defer this.DoFuncLock.Unlock()

		if len(this.DoFuncList) <= 0 {
			return nil
		}

		list := this.DoFuncList
		this.DoFuncList = make([]*timerDo, 0, DO_FUNC_LIST_INIT_CAPACITY)
		return list
	}
	return nil
}

func (this *TimerManager) OnMainLoop() bool {

	list := this.PopDoFuncs()
	if list == nil {
		return true
	}

	for _, do := range list {
		do.Func(do.DoTime, do.Count, do.Para)
	}
	return true
}

func (this *TimerManager) AddTimer(delay int64, delta int64, repeat int64, f TimerFunc, para interface{}) int64 {

	handler := NewTimerHander(this, delay, delta, repeat, f, para)

	// 先添加到缓存列表中
	for this.TempLock.Lock() {
		defer this.TempLock.Unlock()
		this.TempCache.Append(handler)
		this.TimerMap[handler.Id] = handler
	}
	return handler.Id
}

func (this *TimerManager) RemoveTimer(id int64) {

	handler, _ := this.TimerMap[id]
	if handler != nil {
		handler.Deleted = true
		delete(this.TimerMap, id)
	}
}

func (this *TimerManager) newHandlerId() int64 {
	this.NextHandlerId++
	return this.NextHandlerId
}

func (this *TimerManager) addDoTimer(handler *timerHandler) {
	if handler == nil {
		return
	}
	do := newTimeDoFunc(handler.NextDoTime, handler.Count+1, handler.Func, handler.Para)

	for this.DoFuncLock.Lock() {
		defer this.DoFuncLock.Unlock()
		this.DoFuncList = append(this.DoFuncList, do)
		break
	}
}

/******************************************************************************
	timerHandler funcs
******************************************************************************/
func NewTimerHander(mgr *TimerManager, delay int64, delta int64, repeat int64, f TimerFunc, para interface{}) *timerHandler {
	this := new(timerHandler)

	this.Id = mgr.newHandlerId()
	this.CreatTime = _unix_time_ms(nil)
	this.DelayTime = delay
	this.DeltaTime = delta
	this.NextDoTime = this.CreatTime + this.DelayTime

	this.Repeat = repeat
	this.Count = 0 // 从第一次开始计数

	this.Func = f
	this.Para = para

	this.Deleted = false

	this.mgr = mgr
	return this
}

func (this *timerHandler) onUpdate(t int64) {
	if !this.valid() {
		return
	}

	for t > this.NextDoTime {
		this.addDoTimer()

		// 更新数据
		this.Count++
		this.NextDoTime = this.NextDoTime + this.DeltaTime
		// 检查重复情况
		if this.Repeat > 0 {
			if this.Count >= this.Repeat { // 到达重复次数上限, 删除timer
				this.Deleted = true
				break
			}
		}
	}
}

func (this *timerHandler) valid() bool {
	return !this.Deleted
}

func (this *timerHandler) addDoTimer() {
	this.mgr.addDoTimer(this)
}

/******************************************************************************
	TimerDo funcs
******************************************************************************/
func newTimeDoFunc(t int64, count int64, f TimerFunc, para interface{}) *timerDo {
	this := new(timerDo)
	this.DoTime = t
	this.Count = count
	this.Func = f
	this.Para = para
	return this
}
