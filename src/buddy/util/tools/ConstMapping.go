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

type ConstList struct {
	Name   string
	Values []*ConstInfo
}

func NewConstList() *ConstList {
	this := new(ConstList)
	return this
}

type ConstInfo struct {
	Name       string
	KeyType    string
	KeyValue   string
	ValueType  string
	ValueValue string
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
func NewConstMappingFromFile(filename string) *ConstMapping {

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
