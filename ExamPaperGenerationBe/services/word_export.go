package services

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nguyenthenguyen/docx"
)

// 模板文件路径常量
const (
	TemplateTypeTestPaper = 1
	TemplateTypeAnswer    = 2
)

type WordExporter struct {
	data map[string]string
}

func NewWordExporter(data map[string]string) *WordExporter {
	return &WordExporter{data: data}
}

func (we *WordExporter) ExportTestPaper(templateType int) (*os.File, error) {
	// 1. 获取模板路径
	templatePath, err := getTemplatePath(templateType)
	if err != nil {
		return nil, err
	}

	// 2. 加载模板文件
	doc, err := docx.ReadDocxFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	// 3. 替换占位符
	content := doc.Editable()
	for key, value := range we.data {
		placeholder := "${" + key + "}"

		// 处理内容中的图片标记
		if key == "contents" {
			// 使用正则表达式查找所有图片标记
			re := regexp.MustCompile(`\[IMAGE:(.*?)\]`)
			matches := re.FindAllStringSubmatch(value, -1)

			// 替换每个图片标记
			for _, match := range matches {
				imagePath := match[1]
				imagePlaceholder := match[0]

				// 检查图片是否存在
				if _, err := os.Stat(imagePath); os.IsNotExist(err) {
					// 图片不存在，替换为提示文本
					value = strings.Replace(value, imagePlaceholder, "[图片不存在]", -1)
					continue
				}

				// 图片存在，替换为图片占位符
				// 注意：这里我们使用一个特殊的占位符，后面会替换为实际的图片
				value = strings.Replace(value, imagePlaceholder, "[图片]", -1)
			}
		}

		content.Replace(placeholder, value, -1) // -1 表示全部替换
	}

	// 4. 创建临时文件
	tmpFile, err := createTempFile()
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close() // 需要先关闭才能写入

	// 5. 写入临时文件
	if err := content.WriteToFile(tmpPath); err != nil {
		return nil, fmt.Errorf("failed to write document: %w", err)
	}

	// 6. 重新打开文件供调用方使用
	return os.Open(tmpPath)
}

// 获取模板路径（假设模板文件放在项目根目录的 templates 文件夹中）
func getTemplatePath(templateType int) (string, error) {
	var filename string
	switch templateType {
	case TemplateTypeTestPaper:
		filename = "test-paper-template.docx"
	case TemplateTypeAnswer:
		filename = "answer-template.docx"
	default:
		return "", fmt.Errorf("invalid template type: %d", templateType)
	}

	// 假设模板文件路径为 ./templates/[filename]
	path := filepath.Join("templates", filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", path)
	}
	return path, nil
}

// 创建临时文件（带随机后缀）
func createTempFile() (*os.File, error) {
	return os.CreateTemp("", "testpaper_*.docx")
}

//
//// 示例用法
//func main() {
//	data := map[string]string{
//		"title":    "Go语言测试试卷",
//		"question": "1. 什么是Go语言？",
//		"answer":   "Go是一种开源编程语言",
//	}
//
//	exporter := NewWordExporter(data)
//	file, err := exporter.ExportTestPaper(TemplateTypeTestPaper)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//
//	fmt.Printf("生成文件: %s\n", file.Name())
//
//	// 可选：将文件复制到当前目录查看
//	copyToCurrentDir(file.Name())
//}
//
//// 辅助函数：将临时文件复制到当前目录（仅用于测试）
//func copyToCurrentDir(src string) {
//	in, _ := os.Open(src)
//	defer in.Close()
//
//	out, _ := os.Create("output.docx")
//	defer out.Close()
//
//	io.Copy(out, in)
//	fmt.Println("已保存副本到 output.docx")
//}

//package services
//
//import (
//	"bytes"
//	"embed"
//	"fmt"
//	"io"
//	"os"
//	"strings"
//
//	"github.com/unidoc/unioffice/document"
//)
//
////go:embed templates/*.docx
//var templateFS embed.FS
//
//type WordExporter struct {
//	data map[string]string
//}
//
//func NewWordExporter(data map[string]string) *WordExporter {
//	return &WordExporter{data: data}
//}
//
//func (we *WordExporter) ExportTestPaper(templateType int) (*os.File, error) {
//	templatePath, err := getTemplatePath(templateType)
//	if err != nil {
//		return nil, err
//	}
//
//	doc, err := loadTemplate(templatePath)
//	if err != nil {
//		return nil, fmt.Errorf("failed to load template: %w", err)
//	}
//
//	replacePlaceholders(doc, we.data)
//
//	tmpFile, err := createTempFile()
//	if err != nil {
//		return nil, fmt.Errorf("failed to create temp file: %w", err)
//	}
//	defer tmpFile.Close()
//
//	if err := saveDocument(doc, tmpFile); err != nil {
//		return nil, fmt.Errorf("failed to save document: %w", err)
//	}
//
//	return reopenTempFile(tmpFile.Name())
//}
//
//func getTemplatePath(templateType int) (string, error) {
//	switch templateType {
//	case 1:
//		return "templates/test-paper-template.docx", nil
//	case 2:
//		return "templates/answer-template.docx", nil
//	default:
//		return "", fmt.Errorf("invalid template type: %d", templateType)
//	}
//}
//
//func loadTemplate(path string) (*document.Document, error) {
//	file, err := templateFS.Open(path)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	data, err := io.ReadAll(file)
//	if err != nil {
//		return nil, err
//	}
//
//	return document.Read(bytes.NewReader(data), int64(len(data)))
//}
//
//func replacePlaceholders(doc *document.Document, data map[string]string) {
//	for _, p := range doc.Paragraphs() {
//		for _, r := range p.Runs() {
//			original := r.Text()
//			replaced := replaceText(original, data)
//			if replaced != original {
//				r.Clear()
//				r.AddText(replaced)
//			}
//		}
//	}
//}
//
//func replaceText(s string, data map[string]string) string {
//	for k, v := range data {
//		placeholder := "${" + k + "}"
//		s = strings.ReplaceAll(s, placeholder, v)
//	}
//	return s
//}
//
//func createTempFile() (*os.File, error) {
//	return os.CreateTemp("", "testpaper_*.docx")
//}
//
//func saveDocument(doc *document.Document, tmpFile *os.File) error {
//	tmpPath := tmpFile.Name()
//	tmpFile.Close() // Close first to allow SaveToFile to write
//	return doc.SaveToFile(tmpPath)
//}
//
//func reopenTempFile(path string) (*os.File, error) {
//	return os.Open(path)
//}
//
//func main() {
//	// 示例用法
//	data := map[string]string{
//		"title":    "Go语言测试试卷",
//		"question": "1. 什么是Go语言？",
//	}
//
//	exporter := NewWordExporter(data)
//	file, err := exporter.ExportTestPaper(1)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//	fmt.Printf("Generated file: %s\n", file.Name())
//}
