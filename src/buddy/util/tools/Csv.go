package tools

import (
	"encoding/csv"
	"io"
	"os"
)

func CsvLoad(filename string, comma rune) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		ShowError(err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)

	//默认是逗号，也可以自己设置
	if comma != 0 {
		reader.Comma = comma
	}

	list := make([][]string, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			ShowError(err)
			// return nil
		}
		list = append(list, record)
	}
	return list
}
