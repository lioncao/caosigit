package main

import (
	"buddy/util/tools"
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/tealeg_xlsx"
	"os"
	"path/filepath"
	"strings"
)

var list map[string]string
var cfg ConfigMgr

func main() {

	// 读取配置文件
	cfg.Load("main.conf")

	list = make(map[string]string, 0)
	_foreach_file("excels")

	for path, _ := range list {
		_excel_to_csv(path)
	}

}

// 需要处理的文件类型列表
var file_type_list map[string]string = map[string]string{
	"xlsx": "1",
	"xls":  "1",
}

// 对walk中读取到的文件进行处理
func _dofile(path string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}

	name := f.Name()
	strs := strings.Split(name, ".")
	n := len(strs)
	if n < 2 {
		return nil
	}

	fileType := strs[n-1]
	if file_type_list[fileType] == "" {
		return nil
	}

	list[path] = f.Name()
	return nil
}

// 遍历文件列表
func _foreach_file(path string) {
	filepath.Walk(path, _dofile)
}

func _read_excel(filepath string) {

}

// 将excel文件转化为csv文件
func _excel_to_csv(src_excel string) {

	// 将指定excel文件的数据读取出来
	xlFile, err := xlsx.OpenFile(src_excel)
	if err != nil {
		fmt.Println("read file err", src_excel, err.Error())
		return
	}

	names := strings.Split(list[src_excel], ".")

	datas, _ := xlFile.ToSlice()

	// 准备写入到指定文件
	fout, _ := os.Create("csvs/" + names[0] + ".csv")

	buf := new(bytes.Buffer)
	ccc := csv.NewWriter(buf)
	ccc.WriteAll(datas[0])
	ccc.Flush()

	fout.WriteString(buf.String())
	fout.Close()
}
