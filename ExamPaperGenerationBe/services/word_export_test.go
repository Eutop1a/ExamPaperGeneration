package services

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewWordExporter(t *testing.T) {
	exporter := NewWordExporter(map[string]string{
		"title": "Hello, world!",
		"body":  "This is a test.",
	})
	if exporter == nil {
		t.Error("NewWordExporter() should not return nil")
	}
	file, err := exporter.ExportTestPaper(1)
	require.NoError(t, err)
	defer file.Close()
	fmt.Printf("Generated file: %s\n", file.Name())
}

func TestExportTestPaper(t *testing.T) {

	src, err := os.Open("templates/question.xlsx")
	require.NoError(t, err)
	defer src.Close()

	eR := NewExcelReader(src)
	questionBanMap, err := eR.ReadExcel()
	require.NoError(t, err)
	fmt.Println(questionBanMap)
}
