package tools

// 二维参数列表工具， 与配置文件中的ConstMapping格式对应
type ConstMapping struct {
	allConstList map[string]*ConstList
}

func NewConstMapping() *ConstMapping {
	this := new(ConstMapping)
	this.allConstList = make(map[string]*ConstList)
	return this
}

// 单组常数
type ConstList struct {
	Name   string       // 本组的名称
	Values []*ConstInfo // 组中的常数
}

func NewConstList() *ConstList {
	this := new(ConstList)
	return this
}

// 组中的单个常数
type ConstInfo struct {
	Name       string `json:"name"`        // 组名
	KeyType    string `json:"key_type"`    // key类型
	KeyValue   string `json:"key_value"`   // key值的字符串形式
	ValueType  string `json:"value_type"`  // value类型
	ValueValue string `json:"value_value"` // value对应的字符串形式
}

func NewConstInfo() *ConstInfo {
	this := new(ConstInfo)
	return this
}

func (this *ConstInfo) Init(cd *CommonData) {
	this.Name = cd.S("name")
	this.KeyType = cd.S("key_type")
	this.KeyValue = cd.S("key_value")
	this.ValueType = cd.S("value_type")
	this.ValueValue = cd.S("value_value")
}

func NewConstInfoFromCommonData(cd *CommonData) *ConstInfo {
	this := NewConstInfo()
	this.Init(cd)
	return this
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// 从ConstMapping配置文件中构建一个ConstMapping对象
func NewConstMappingFromJsonObjFile(filename string) *ConstMapping {
	this := NewConstMapping()

	// 从 json_obj文件中读出数据
	datas := make([]*ConstInfo, 0)
	LoadJsonFile(filename, datas)
	cnt := len(datas)

	// 数据入表
	var (
		info *ConstInfo
		list *ConstList
	)

	for i := 0; i < cnt; i++ {
		info = datas[i]
		list = this.allConstList[info.Name]
		if list == nil {
			list = NewConstList()
			this.allConstList[info.Name] = list
		}
		list.Values = append(list.Values, info)
	}

	return this
}

// 从ConstMapping配置文件中构建一个ConstMapping对象
func NewConstMappingFromCommonJsonFile(filename string) *ConstMapping {

	cdList := NewCommonJsonDataFromFile(filename).ToCommDataList()
	if cdList == nil {
		return nil
	}
	count := len(cdList)

	this := NewConstMapping()

	var (
		list *ConstList
		cd   *CommonData
		info *ConstInfo
	)

	for i := 0; i < count; i++ {
		cd = cdList[i]
		info = NewConstInfoFromCommonData(cd)
		list = this.allConstList[info.Name]
		if list == nil {
			list = NewConstList()
			this.allConstList[info.Name] = list
		}
		list.Values = append(list.Values, info)
	}
	return this
}
func (this *ConstMapping) GetConstList(name string) *ConstList {
	list, ok := this.allConstList[name]
	if ok {
		return list
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func (this *ConstList) ToList() []string {
	count := len(this.Values)
	list := make([]string, count)
	var (
		info *ConstInfo
	)
	for i := 0; i < count; i++ {
		info = this.Values[i]
		list[i] = info.ValueValue
	}
	return list
}

func (this *ConstList) ToKeyList() []string {
	count := len(this.Values)
	list := make([]string, count)
	var (
		info *ConstInfo
	)
	for i := 0; i < count; i++ {
		info = this.Values[i]
		list[i] = info.KeyValue
	}
	return list
}

func (this *ConstList) ToMap() map[string]string {
	count := len(this.Values)
	valuemap := make(map[string]string)
	var (
		info     *ConstInfo
		valueOrg string
		ok       bool
	)
	for i := 0; i < count; i++ {
		info = this.Values[i]
		valueOrg, ok = valuemap[info.KeyValue]
		if ok { // 取到了同名值
			ShowWarnning("ConstListToMap  key redefined", info.Name, info.KeyValue, valueOrg, info.ValueValue)
		}
		valuemap[info.KeyValue] = info.ValueValue
	}
	return valuemap
}
