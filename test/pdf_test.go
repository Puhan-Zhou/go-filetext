package test

import (
	"testing"

	"github.com/Puhan-Zhou/go-filetext/extractor"
)

func TestPDFExtraction(t *testing.T) {
	pdfExtractor := &extractor.PDFExtractor{}
	options := extractor.DefaultExtractOptions()

	// Test PDF extraction from file
	result, err := pdfExtractor.ExtractFromFile("testdata/sample.pdf", options)
	if err != nil {
		t.Fatalf("PDF extraction failed: %v", err)
	}

	if result.Text == "" {
		t.Error("Expected non-empty text from PDF")
	}

	if result.FileType != "pdf" {
		t.Errorf("Expected file type 'pdf', got '%s'", result.FileType)
	}

	if result.Metadata["page_count"] == nil {
		t.Error("Expected page_count metadata")
	}

	t.Logf("PDF extraction successful. Text length: %d", len(result.Text))
	t.Logf("Metadata: %+v", result.Metadata)
}

func TestPDFExtractionWithOptions(t *testing.T) {
	pdfExtractor := &extractor.PDFExtractor{}
	options := extractor.DefaultExtractOptions()
	options.MaxFileSize = 1 // Very small limit to test size restriction

	// Test with size limit
	_, err := pdfExtractor.ExtractFromFile("testdata/sample.pdf", options)
	if err == nil {
		t.Error("Expected error due to file size limit")
	}
}
