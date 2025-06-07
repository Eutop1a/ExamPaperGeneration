package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/carmel/gooxml/common"
	"github.com/carmel/gooxml/document"
	measure "github.com/carmel/gooxml/measurement"
	"github.com/carmel/gooxml/schema/soo/wml"
)

// WordExporterGooxml 使用 gooxml 库导出 Word 文档
type WordExporterGooxml struct {
	data map[string]string
}

// NewWordExporterGooxml 创建一个新的 WordExporterGooxml 实例
func NewWordExporterGooxml(data map[string]string) *WordExporterGooxml {
	return &WordExporterGooxml{data: data}
}

// ExportTestPaper 导出试卷
func (we *WordExporterGooxml) ExportTestPaper(templateType int) (*os.File, error) {
	// 创建新文档
	doc := document.New()

	// 添加试卷基本信息
	// 添加试题编号 - 宋体4号
	numberPara := doc.AddParagraph()
	numberPara.Properties().SetAlignment(wml.ST_JcLeft)
	numberRun := numberPara.AddRun()
	numberRun.Properties().SetSize(14 * measure.Point) // 4号字体约为14pt
	numberRun.Properties().SetFontFamily("宋体")
	numberRun.AddText("试题编号：")

	// 添加学校、学年、学期信息 - 黑体2号
	schoolPara := doc.AddParagraph()
	schoolPara.Properties().SetAlignment(wml.ST_JcCenter)
	schoolRun := schoolPara.AddRun()
	schoolRun.Properties().SetSize(22 * measure.Point) // 2号字体约为22pt
	schoolRun.Properties().SetFontFamily("黑体")
	schoolRun.AddText("重庆邮电大学XXXX学年第X学期（试卷）")

	// 添加课程信息 - 黑体三号
	coursePara := doc.AddParagraph()
	coursePara.Properties().SetAlignment(wml.ST_JcCenter)
	courseRun := coursePara.AddRun()
	courseRun.Properties().SetSize(16 * measure.Point) // 三号字体约为16pt
	courseRun.Properties().SetFontFamily("黑体")
	courseRun.AddText("XXXX课程（期末/期中）（A/B卷）（开卷/闭卷/其他）")

	// 添加分隔线
	doc.AddParagraph()

	// 添加标题
	titlePara := doc.AddParagraph()
	titlePara.Properties().SetAlignment(wml.ST_JcCenter)
	titleRun := titlePara.AddRun()
	titleRun.Properties().SetBold(true)
	titleRun.Properties().SetSize(12 * measure.Point) // 5号字体约为10.5pt，这里使用12pt
	titleRun.AddText(we.data["test_paper_name"])

	// 添加表头信息
	headerTable := doc.AddTable()
	headerTable.Properties().SetAlignment(wml.ST_JcTableCenter)
	headerTable.Properties().SetWidth(12 * measure.Centimeter)

	// 添加表格行
	row := headerTable.AddRow()

	// 添加总分单元格
	cell := row.AddCell()
	cellPara := cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellRun := cellPara.AddRun()
	cellRun.Properties().SetSize(10.5 * measure.Point) // 5号字体
	cellRun.AddText(fmt.Sprintf("总分：%s分", we.data["total_score"]))

	// 添加题目数量单元格
	cell = row.AddCell()
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellRun = cellPara.AddRun()
	cellRun.Properties().SetSize(10.5 * measure.Point) // 5号字体
	cellRun.AddText(fmt.Sprintf("题目数量：%s题", we.data["total_count"]))

	// 添加分隔线
	doc.AddParagraph()

	// 处理题目内容
	contents := we.data["contents"]
	questions := strings.Split(contents, "[QUESTION_END]")
	re := regexp.MustCompile(`\[IMAGE:(.*?)\]`)
	sectionTitleRe := regexp.MustCompile(`\[SECTION_TITLE\](.*?)\[/SECTION_TITLE\]`)

	// 用于跟踪当前是否在简答题部分
	isShortAnswerSection := false
	// 用于跟踪当前题目是否已经添加了答题空间
	addedAnswerSpace := false

	for _, question := range questions {
		question = strings.TrimSpace(question)
		if question == "" {
			continue
		}

		// 检查是否包含题型描述
		if strings.Contains(question, "[SECTION_TITLE]") {
			// 提取标题内容
			matches := sectionTitleRe.FindStringSubmatch(question)
			if len(matches) > 1 {
				// 添加前空行
				doc.AddParagraph()

				// 添加标题
				para := doc.AddParagraph()
				para.Properties().SetAlignment(wml.ST_JcLeft)
				run := para.AddRun()
				run.Properties().SetBold(true)                 // 大题标题加粗
				run.Properties().SetSize(10.5 * measure.Point) // 5号字体
				run.AddText(matches[1])

				// 添加后空行
				doc.AddParagraph()

				// 检查是否是简答题部分
				if strings.Contains(question, "简答题") {
					isShortAnswerSection = true
				} else {
					isShortAnswerSection = false
				}

				// 移除题型描述部分，只保留题目内容
				question = strings.TrimSpace(sectionTitleRe.ReplaceAllString(question, ""))
				if question == "" {
					continue
				}
			}
		}

		// 检查是否包含图片标记
		if strings.Contains(question, "[IMAGE:") {
			// 提取图片路径
			var imgPath string
			matches := re.FindStringSubmatch(question)
			if len(matches) > 1 {
				imgPath = matches[1]
			}

			// 创建一个段落来包含题目文本和图片
			para := doc.AddParagraph()
			para.Properties().SetAlignment(wml.ST_JcLeft)
			run := para.AddRun()
			run.Properties().SetSize(10.5 * measure.Point) // 5号字体
			// 确保小题不加粗
			run.Properties().SetBold(false)

			// 添加题目文本（移除图片标记）
			text := strings.TrimSpace(re.ReplaceAllString(question, ""))
			run.AddText(text)

			// 读取并添加图片
			imgData, err := ioutil.ReadFile(imgPath)
			if err != nil {
				log.Printf("Error reading image file %s: %v", imgPath, err)
				continue
			}

			img, err := common.ImageFromBytes(imgData)
			if err != nil {
				log.Printf("Error creating image from bytes: %v", err)
				continue
			}

			imgRef, err := doc.AddImage(img)
			if err != nil {
				log.Printf("Error adding image to document: %v", err)
				continue
			}

			// 在同一个段落中添加图片
			inline, err := run.AddDrawingInline(imgRef)
			if err != nil {
				log.Printf("Error adding inline image: %v", err)
				continue
			}

			inline.SetSize(2*measure.Inch, 1.5*measure.Inch)
		} else {
			// 按行分割文本
			lines := strings.Split(question, "\n")
			// 处理普通题目
			for i, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				// 创建段落
				para := doc.AddParagraph()
				para.Properties().SetAlignment(wml.ST_JcLeft)
				run := para.AddRun()
				run.Properties().SetSize(10.5 * measure.Point) // 5号字体
				// 确保小题不加粗
				run.Properties().SetBold(false)

				// 处理选项的缩进
				if strings.HasPrefix(line, "①") || strings.HasPrefix(line, "②") ||
					strings.HasPrefix(line, "③") || strings.HasPrefix(line, "④") ||
					strings.HasPrefix(line, "A.") || strings.HasPrefix(line, "B.") ||
					strings.HasPrefix(line, "C.") || strings.HasPrefix(line, "D.") {
					// 选项前添加制表符
					run.AddText("\t" + line)
				} else if strings.HasPrefix(line, "（") && strings.Contains(line, "）") {
					// 小题编号前添加制表符
					run.AddText("\t" + line)
				} else {
					run.AddText(line)
				}

				// 如果是简答题的最后一行，添加答题空间
				if isShortAnswerSection && i == len(lines)-1 && !addedAnswerSpace {
					// 添加多个空行作为答题空间
					for j := 0; j < 8; j++ {
						spacePara := doc.AddParagraph()
						spacePara.Properties().SetAlignment(wml.ST_JcLeft)
						spaceRun := spacePara.AddRun()
						spaceRun.Properties().SetSize(10.5 * measure.Point) // 5号字体
						spaceRun.AddText("")
					}
					addedAnswerSpace = true
				} else if i == len(lines)-1 {
					// 重置状态
					addedAnswerSpace = false
				}
			}
		}
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "testpaper_*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	// 保存文档
	if err := doc.SaveToFile(tmpFile.Name()); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	// 重新打开文件供调用方使用
	return os.Open(tmpFile.Name())
}

// ExportAnswer 导出答案
func (we *WordExporterGooxml) ExportAnswer(templateType int) (*os.File, error) {
	// 创建新文档
	doc := document.New()

	// 添加标题
	titlePara := doc.AddParagraph()
	titlePara.Properties().SetAlignment(wml.ST_JcCenter)
	titleRun := titlePara.AddRun()
	titleRun.Properties().SetBold(true)
	titleRun.Properties().SetSize(12 * measure.Point) // 5号字体约为10.5pt，这里使用12pt
	titleRun.AddText("答案")

	// 添加表头信息
	headerTable := doc.AddTable()
	headerTable.Properties().SetAlignment(wml.ST_JcTableCenter)
	headerTable.Properties().SetWidth(12 * measure.Centimeter)

	// 添加表格行
	row := headerTable.AddRow()

	// 添加总分单元格
	cell := row.AddCell()
	cellPara := cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellRun := cellPara.AddRun()
	cellRun.Properties().SetSize(10.5 * measure.Point) // 5号字体
	cellRun.AddText(fmt.Sprintf("总分：%s分", we.data["total_score"]))

	// 添加题目数量单元格
	cell = row.AddCell()
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellRun = cellPara.AddRun()
	cellRun.Properties().SetSize(10.5 * measure.Point) // 5号字体
	cellRun.AddText(fmt.Sprintf("题目数量：%s题", we.data["total_count"]))

	// 添加分隔线
	doc.AddParagraph()

	// 处理题目内容
	contents := we.data["contents"]
	questions := strings.Split(contents, "[QUESTION_END]")
	re := regexp.MustCompile(`\[IMAGE:(.*?)\]`)

	for _, question := range questions {
		question = strings.TrimSpace(question)
		if question == "" {
			continue
		}

		// 检查是否包含图片标记
		if strings.Contains(question, "[IMAGE:") {
			// 提取图片路径
			var imgPath string
			matches := re.FindStringSubmatch(question)
			if len(matches) > 1 {
				imgPath = matches[1]
			}

			// 创建一个段落来包含题目文本和图片
			para := doc.AddParagraph()
			para.Properties().SetAlignment(wml.ST_JcLeft)
			run := para.AddRun()
			run.Properties().SetSize(10.5 * measure.Point) // 5号字体

			// 添加题目文本（移除图片标记）
			text := strings.TrimSpace(re.ReplaceAllString(question, ""))
			run.AddText(text)

			// 读取并添加图片
			imgData, err := ioutil.ReadFile(imgPath)
			if err != nil {
				log.Printf("Error reading image file %s: %v", imgPath, err)
				continue
			}

			img, err := common.ImageFromBytes(imgData)
			if err != nil {
				log.Printf("Error creating image from bytes: %v", err)
				continue
			}

			imgRef, err := doc.AddImage(img)
			if err != nil {
				log.Printf("Error adding image to document: %v", err)
				continue
			}

			// 在同一个段落中添加图片
			inline, err := run.AddDrawingInline(imgRef)
			if err != nil {
				log.Printf("Error adding inline image: %v", err)
				continue
			}

			inline.SetSize(2*measure.Inch, 1.5*measure.Inch)
		} else {
			// 按行分割文本
			lines := strings.Split(question, "\n")

			// 检查是否是大题标题
			isSectionTitle := false
			for _, line := range lines {
				if strings.HasPrefix(line, "一、") || strings.HasPrefix(line, "二、") ||
					strings.HasPrefix(line, "三、") || strings.HasPrefix(line, "四、") {
					isSectionTitle = true
					break
				}
			}

			// 如果是大题标题，添加前后空行
			if isSectionTitle {
				// 添加前空行
				doc.AddParagraph()

				// 添加标题
				for _, line := range lines {
					para := doc.AddParagraph()
					para.Properties().SetAlignment(wml.ST_JcLeft)
					run := para.AddRun()
					run.Properties().SetBold(true)                 // 大题标题加粗
					run.Properties().SetSize(10.5 * measure.Point) // 5号字体
					run.AddText(line)
				}

				// 添加后空行
				doc.AddParagraph()

				continue
			}

			// 处理普通题目
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				// 创建段落
				para := doc.AddParagraph()
				para.Properties().SetAlignment(wml.ST_JcLeft)
				run := para.AddRun()
				run.Properties().SetSize(10.5 * measure.Point) // 5号字体

				// 处理选项的缩进
				if strings.HasPrefix(line, "①") || strings.HasPrefix(line, "②") ||
					strings.HasPrefix(line, "③") || strings.HasPrefix(line, "④") {
					// 选项前添加制表符
					run.AddText("\t" + line)
				} else if strings.HasPrefix(line, "（") && strings.Contains(line, "）") {
					// 小题编号前添加制表符
					run.AddText("\t" + line)
				} else {
					run.AddText(line)
				}
			}
		}
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "answer_*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	// 保存文档
	if err := doc.SaveToFile(tmpFile.Name()); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	// 重新打开文件供调用方使用
	return os.Open(tmpFile.Name())
}
