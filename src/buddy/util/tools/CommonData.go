package tools

import (
	"encoding/json"
	// "errors"
	"fmt"
	"os"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// 通用db相关数据结构
////////////////////////////////////////////////////////////////////////////////
type CommonDataPool struct {
	datas map[string]*CommonData
}

func (this *CommonDataPool) Init() {
	this.datas = make(map[string]*CommonData)
}

func (this *CommonDataPool) Put(key string, data *CommonData) error {
	if key == "" || data == nil {
		return Error("put common data err key=%s,data=%v", key, data)
	}
	this.datas[key] = data
	return nil
}

func (this *CommonDataPool) Get(key string) *CommonData {
	if key == "" {
		return nil
	}
	v, ok := this.datas[key]
	if !ok {
		return nil
	}
	return v
}

// 统一数据格式, 对应data.json
type CommonData struct {
	FileName string     `json:"-"`      // 原始文件的名字
	JsonStr  string     `json:"-"`      // 原始的json字符串
	Fields   []string   `json:"fields"` // 数据名称列表
	Types    []string   `json:"types"`  // 数据类型列表
	Values   [][]string `json:"values"` // 数据表
	ValueBuf [][]byte   `json:"-"`      // 每个value对应的json字符串(相当于预先打包好的单条数据的jsonStr)

	// 辅助数据
	DataCount    int            `json:"-"` // 数据总条数
	FieldNameMap map[string]int `json:"-"` // 数据名称到数据下标的映射
	FieldCount   int            `json:"-"` // 一条数据的数据个数

}

func NewCommonData() *CommonData {
	this := new(CommonData)
	// this.FieldNameMap = new(map[string]int)
	return this
}

// 从json格式中读取统一数据
func (this *CommonData) DecodeJsonFile(filename string) error {
	var e error

	f, e := os.Open(filename)
	if e != nil {
		return e
	}
	defer f.Close()

	var buf, jsonBuffer []byte

	bufLen := 1024
	buf = make([]byte, bufLen) // read buf

	var n, count int
	count = 0
	for {
		n, e = f.Read(buf)
		if e != nil || n <= 0 {
			break
		}

		if jsonBuffer == nil { // 首次读取
			if n < bufLen {
				jsonBuffer = buf[0:n]
				break
			}
			jsonBuffer = make([]byte, bufLen<<2)
		}

		jsonBuffer = append(jsonBuffer[0:count], buf[:n]...)
		count += n
	}

	// check uft8 bom
	if CheckUTF8_BOM(jsonBuffer) {
		jsonBuffer = jsonBuffer[UTF8_BOM_LEN:]
	}

	e = json.Unmarshal(jsonBuffer, this)
	if e != nil {
		return e
	}
	this.FileName = filename
	this.JsonStr = string(jsonBuffer)

	// field 总数
	if this.Fields != nil {
		//  field 整理
		this.FieldCount = len(this.Fields)
		this.FieldNameMap = make(map[string]int, this.FieldCount)
		for i := 0; i < this.FieldCount; i++ {
			this.FieldNameMap[this.Fields[i]] = i
		}
	} else {
		this.FieldCount = 0
	}

	// 数据条数
	if this.Values != nil {
		this.DataCount = len(this.Values)
	} else {
		this.DataCount = 0
	}

	return nil
}

func (this *CommonData) getValue(index int, fieldName string) string {
	var (
		fieldIndex int
		ok         bool
	)

	if index < 0 || index > this.DataCount {
		ShowError("CommonData.getValue() invalid index , fileName=", Color(CL_YELLOW, this.FileName),
			", index=", index, ", DataCount=", this.DataCount)
		return ""
	}

	fieldIndex, ok = this.FieldNameMap[fieldName]
	if !ok {
		ShowError("CommonData.getValue() invalid fieldName , fileName=", Color(CL_YELLOW, this.FileName),
			" ,fieldName=", fieldName)
		return ""
	}

	return this.Values[index][fieldIndex]
}

// 从数据集中解析出一个 string
func (this *CommonData) ParseString(index int, fieldName string) string {
	return this.getValue(index, fieldName)
}

// 从数据集中解析出一个int64
func (this *CommonData) ParseInt64(index int, fieldName string) int64 {
	var (
		value string
	)

	value = this.getValue(index, fieldName)
	if value == "" {
		ShowWarnning("CommonData.ParseInt64() get empty value , fileName=", Color(CL_YELLOW, this.FileName),
			",index=", index, ",fieldName=", fieldName, ",ID=", this.getValue(index, "ID"))
		return 0
	}

	i64, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		ShowError("CommonData.ParseInt64() get value faild, fileName=", Color(CL_YELLOW, this.FileName),
			",index=", index, ",fieldName=", fieldName, "value=", value, ",ID=", this.getValue(index, "ID"))
		return 0
	}

	return i64
}

// 从数据集中解析出一个 float64
func (this *CommonData) ParseFloat64(index int, fieldName string) float64 {
	var (
		value string
	)

	value = this.getValue(index, fieldName)
	if value == "" {
		ShowWarnning("CommonData.ParseFloat64() get empty value , fileName=", Color(CL_YELLOW, this.FileName),
			",index=", index, ",fieldName=", fieldName, ",ID=", this.getValue(index, "ID"))
		return 0.0
	}

	f64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		ShowError("CommonData.ParseFloat64() get value faild, fileName=", Color(CL_YELLOW, this.FileName),
			",index=", index, ",fieldName=", fieldName, "value=", value, ",ID=", this.getValue(index, "ID"))
		return 0.0
	}

	return f64
}

// 从数据集中解析出一个 bool
func (this *CommonData) ParseBool(index int, fieldName string) bool {
	var (
		value string
	)

	value = this.getValue(index, fieldName)
	if value == "" {
		ShowWarnning("CommonData.ParseBool() get empty value , fileName=", Color(CL_YELLOW, this.FileName),
			",index=", index, ",fieldName=", fieldName, ",ID=", this.getValue(index, "ID"))
		return false
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		ShowError("CommonData.ParseFloat64() get value faild, fileName=", Color(CL_YELLOW, this.FileName),
			",index=", index, ",fieldName=", fieldName, "value=", value, ",ID=", this.getValue(index, "ID"))
		return false
	}

	return b
}

func (this *CommonData) Print() {
	ShowDebug(fmt.Sprintf("\n\n-----start print file \"%s\"------------------------", this.FileName))
	if this.Fields != nil {
		ShowDebug("fields <", len(this.Fields), "> data")
		for k, v := range this.Fields {
			ShowDebug("\t", k, v, "\t")
		}
	}

	if this.Types != nil {
		ShowDebug("types <", len(this.Types), "> data")
		for k, v := range this.Types {
			ShowDebug("\t", k, v, "\t")
		}
	}

	if this.Values != nil {
		ShowDebug("values <", len(this.Values), "> data")
		for k, v := range this.Values {
			ShowDebug("\tvalue", k, " <", len(v), "> data")
			for x, y := range v {
				ShowDebug("\t\tvalues", k, x, y, "\t")
			}
		}
	}
	ShowDebug(fmt.Sprintf("\n-----end print file \"%s\"------------------------\n\n", this.FileName))

}
