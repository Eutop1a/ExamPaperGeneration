package services

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"strconv"
	"testing"
)

func TestReadExcel(t *testing.T) {
	filePath := "../resources/QuestionBank.xlsx"
	file, err := os.Open(filePath)
	require.NoError(t, err)
	ert := NewExcelReaderTest(file)
	res, err := ert.ReadExcel()
	require.NoError(t, err)
	spew.Dump(res)
}

// ExcelReader 结构体用于读取 Excel 文件
type ExcelReaderTest struct {
	inputStream io.Reader
}

// NewExcelReader 创建一个新的 ExcelReader 实例
func NewExcelReaderTest(inputStream io.Reader) *ExcelReaderTest {
	return &ExcelReaderTest{
		inputStream: inputStream,
	}
}

// ReadExcel 读取 Excel 文件并返回数据
func (er *ExcelReaderTest) ReadExcel() ([]map[string]interface{}, error) {
	f, err := excelize.OpenReader(er.inputStream)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// 获取第一个工作表的名称
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	// 从第二行开始读取数据
	for idx, row := range rows[1:] {
		r := make(map[string]interface{})
		if len(row) > 1 {
			r["topic"] = row[0]
		}
		if len(row) > 2 {
			r["topic_material_id"] = row[1]
		}
		if len(row) > 3 {
			r["answer"] = row[2]
		}
		if len(row) > 4 {
			r["topic_type"] = row[3]
		}
		if len(row) > 5 {
			score, err := strconv.ParseFloat(row[4], 64)
			if err == nil {
				r["score"] = score
			}
		}
		if len(row) > 6 {
			difficulty, err := strconv.ParseInt(row[5], 10, 64)
			if err == nil {
				r["difficulty"] = difficulty
			}
		}
		if len(row) > 7 {
			r["chapter_1"] = row[6]
		}
		if len(row) > 8 {
			r["chapter_2"] = row[7]
		}
		if len(row) > 9 {
			r["label_1"] = row[8]
		}
		if len(row) > 10 {
			r["label_2"] = row[9]
		}
		if len(row) > 11 {
			r["update_time"] = row[10]
		}
		if len(row) >= 12 {
			r["image"] = row[11]
			// 获取所有图片
			pics, err := f.GetPictures(sheetName, fmt.Sprintf("L%d", idx+2))
			if err != nil {
				fmt.Println("Error getting images:", err)
				return nil, err
			}
			for i, pic := range pics {
				name := fmt.Sprintf("../resources/images/image%d%s", idx+i+1, pic.Extension)
				if err := os.WriteFile(name, pic.File, 0644); err != nil {
					fmt.Println(err)
					return nil, err
				}
			}
		}
		result = append(result, r)
	}
	return result, nil
}
