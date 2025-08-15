package test

import (
	"testing"

	"github.com/Puhan-Zhou/go-filetext/pkg/extractor"
)

func TestDOCXExtraction(t *testing.T) {
	docxExtractor := &extractor.DOCXExtractor{}
	options := extractor.DefaultExtractOptions()

	// Test DOCX extraction from file
	result, err := docxExtractor.ExtractFromFile("testdata/sample.docx", options)
	if err != nil {
		t.Fatalf("DOCX extraction failed: %v", err)
	}

	if result.Text == "" {
		t.Error("Expected non-empty text from DOCX")
	}

	if result.FileType != "docx" {
		t.Errorf("Expected file type 'docx', got '%s'", result.FileType)
	}

	if result.Metadata["paragraphs"] == nil {
		t.Error("Expected paragraphs metadata")
	}

	if result.Metadata["characters"] == nil {
		t.Error("Expected characters metadata")
	}

	t.Logf("DOCX extraction successful. Text length: %d", len(result.Text))
	t.Logf("Metadata: %+v", result.Metadata)
}

func TestDOCXExtractionWithOptions(t *testing.T) {
	docxExtractor := &extractor.DOCXExtractor{}
	options := extractor.DefaultExtractOptions()
	options.MaxFileSize = 1 // Very small limit to test size restriction

	// Test with size limit
	_, err := docxExtractor.ExtractFromFile("testdata/sample.docx", options)
	if err == nil {
		t.Error("Expected error due to file size limit")
	}
}
