package services

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
)

// 系统路径配置
const (
	WindowsBasePath = "D:/graduation/resources/"
	LinuxBasePath   = "/graduation/resources/"
	ImageSubPath    = "images/"
	TableSubPath    = "tables/"
)

// ExcelReader 用于读取Excel文件数据
type ExcelReader struct {
	inputStream io.Reader
}

// NewExcelReader 创建ExcelReader实例
func NewExcelReader(inputStream io.Reader) *ExcelReader {
	return &ExcelReader{
		inputStream: inputStream,
	}
}

// ReadExcel 读取Excel文件并返回结构化数据
func (er *ExcelReader) ReadExcel() ([]map[string]interface{}, error) {
	f, err := excelize.OpenReader(er.inputStream)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("excel文件中没有工作表")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表数据失败: %w", err)
	}

	if len(rows) <= 1 {
		return nil, nil // 没有数据行
	}

	return er.processRows(rows[1:]), nil
}

// processRows 处理行数据并转换为map切片
func (er *ExcelReader) processRows(rows [][]string) []map[string]interface{} {
	var result []map[string]interface{}

	for _, row := range rows {
		if len(row) == 0 {
			continue
		}

		record := make(map[string]interface{})

		// 使用列索引常量提高可读性
		const (
			topicCol      = 0
			materialCol   = 1
			answerCol     = 2
			typeCol       = 3
			scoreCol      = 4
			difficultyCol = 5
			chapter1Col   = 6
			chapter2Col   = 7
			label1Col     = 8
			label2Col     = 9
			updateTimeCol = 10
		)

		// 处理各列数据
		if len(row) > topicCol {
			record["topic"] = row[topicCol]
		}
		if len(row) > materialCol {
			record["topic_material_id"] = row[materialCol]
		}
		if len(row) > answerCol {
			record["answer"] = row[answerCol]
		}
		if len(row) > typeCol {
			record["topic_type"] = row[typeCol]
		}
		if len(row) > scoreCol {
			if score, err := strconv.ParseFloat(row[scoreCol], 64); err == nil {
				record["score"] = score
			}
		}
		if len(row) > difficultyCol {
			if difficulty, err := strconv.ParseInt(row[difficultyCol], 10, 64); err == nil {
				record["difficulty"] = difficulty
			}
		}
		if len(row) > chapter1Col {
			record["chapter_1"] = row[chapter1Col]
		}
		if len(row) > chapter2Col {
			record["chapter_2"] = row[chapter2Col]
		}
		if len(row) > label1Col {
			record["label_1"] = row[label1Col]
		}
		if len(row) > label2Col {
			record["label_2"] = row[label2Col]
		}
		if len(row) > updateTimeCol {
			record["update_time"] = row[updateTimeCol]
		}

		result = append(result, record)
	}

	return result
}
