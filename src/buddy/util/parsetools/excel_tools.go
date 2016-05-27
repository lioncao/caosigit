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

var (
	cfg tools.ConfigMgr // 配置管理器

	// 源目录和输出目录
	src_path string
	tag_path string

	// 需要处理的文件类型列表
	file_type_list map[string]string = map[string]string{
		"xlsx": "xlsx",
		"xls":  "xls",
	}

	// 解析的表格数量
	parse_len int64 = 1

	// 用于缓存文件列表
	list map[string]string
)

func main() {

	// 读取配置文件
	_load_cfg()

	list = make(map[string]string, 0)
	_foreach_file(src_path)

	for path, _ := range list {
		_excel_to_csv(path)
	}

}

func _load_cfg() {
	cfg.Load("main.conf")
	cfg.Print()

	src_path = cfg.EnsureString("global", "src_path", "excels")
	tag_path = cfg.EnsureString("global", "tag_path", "csvs")

	types := tools.SliptStr(cfg.EnsureString("global", "src_file_types", "xlsx,xls"), ",")
	if len(types) > 0 {
		file_type_list = make(map[string]string, 0)
		for _, t := range types {
			file_type_list[t] = t
		}
	}

	parse_len = cfg.EnsureInt("global", "parse_len", 1)

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

// 将excel文件转化为csv文件
func _excel_to_csv(src_excel string) {

	// 将指定excel文件的数据读取出来
	xlFile, err := xlsx.OpenFile(src_excel)
	if err != nil {
		fmt.Println("read file err", src_excel, err.Error())
		return
	}

	// names := strings.Split(list[src_excel], ".")

	datas, _ := xlFile.ToSlice()

	sheetCnt := len(datas)
	cnt := int(parse_len)
	if cnt <= 0 || cnt > sheetCnt {
		cnt = sheetCnt
	}

	var (
		foutName string
	)
	for i := 0; i < cnt; i++ {
		foutName = strings.TrimSpace(xlFile.Sheets[i].Name)
		if strings.HasPrefix(foutName, "@") {
			tools.ShowWarnning("skip sheet ", src_excel, foutName)
			continue
		}

		tools.ShowInfo("trance sheet", src_excel, foutName)

		if i == 0 || true {
			foutName = fmt.Sprintf("%s/%s.csv", tag_path, foutName)
		} else {
			foutName = fmt.Sprintf("%s/%s_%d.csv", tag_path, foutName, i)
		}

		_, e := os.Stat(foutName)
		if e == nil {
			tools.Error("output file exist", foutName)
			continue
		}

		// 准备写入到指定文件
		fout, err := os.Create(foutName)
		defer fout.Close()

		if err != nil {
			tools.ShowError(err.Error())
		}

		buf := new(bytes.Buffer)
		ccc := csv.NewWriter(buf)
		ccc.WriteAll(datas[i])
		ccc.Flush()

		fout.WriteString(buf.String())
	}

}
