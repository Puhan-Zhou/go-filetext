package test

import (
	"testing"

	"github.com/Puhan-Zhou/go-file-plain-text/pkg/extractor"
)

func TestPPTXExtraction(t *testing.T) {
	pptxExtractor := &extractor.PPTXExtractor{}
	options := extractor.DefaultExtractOptions()

	// Test PPTX extraction from file
	result, err := pptxExtractor.ExtractFromFile("testdata/sample.pptx", options)
	if err != nil {
		t.Fatalf("PPTX extraction failed: %v", err)
	}

	if result.Text == "" {
		t.Error("Expected non-empty text")
	}

	if result.FileType != "pptx" {
		t.Errorf("Expected file type 'pptx', got '%s'", result.FileType)
	}

	// Check metadata
	if result.Metadata == nil {
		t.Error("Expected metadata to be present")
	}

	// Check that we have slide count
	if slides, ok := result.Metadata["slides"]; !ok || slides == "0" {
		t.Error("Expected slides metadata to be present and non-zero")
	}

	// Check that we have character count
	if chars, ok := result.Metadata["characters"]; !ok || chars == "0" {
		t.Error("Expected characters metadata to be present and non-zero")
	}
}
