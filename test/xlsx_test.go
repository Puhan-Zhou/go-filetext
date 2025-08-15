package test

import (
	"testing"

	"github.com/Puhan-Zhou/go-filetext/extractor"
)

func TestXLSXExtraction(t *testing.T) {
	xlsxExtractor := &extractor.XLSXExtractor{}
	options := extractor.DefaultExtractOptions()

	// Test XLSX extraction from file
	result, err := xlsxExtractor.ExtractFromFile("testdata/sample.xlsx", options)
	if err != nil {
		t.Fatalf("XLSX extraction failed: %v", err)
	}

	if result.Text == "" {
		t.Error("Expected non-empty text from XLSX")
	}

	if result.FileType != "xlsx" {
		t.Errorf("Expected file type 'xlsx', got '%s'", result.FileType)
	}

	if result.Metadata["sheets"] == nil {
		t.Error("Expected sheets metadata")
	}

	if result.Metadata["rows"] == nil {
		t.Error("Expected rows metadata")
	}

	if result.Metadata["cells"] == nil {
		t.Error("Expected cells metadata")
	}

	t.Logf("XLSX extraction successful. Text length: %d", len(result.Text))
	t.Logf("Metadata: %+v", result.Metadata)
}

func TestXLSXExtractionWithOptions(t *testing.T) {
	xlsxExtractor := &extractor.XLSXExtractor{}
	options := extractor.DefaultExtractOptions()
	options.MaxFileSize = 1 // Very small limit to test size restriction

	// Test with size limit
	_, err := xlsxExtractor.ExtractFromFile("testdata/sample.xlsx", options)
	if err == nil {
		t.Error("Expected error due to file size limit")
	}
}
