package test

import (
	"path/filepath"
	"testing"

	"github.com/Puhan-Zhou/go-file-plain-text/pkg/extractor"
)

func TestImageExtraction(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
		expected string
	}{
		{
			name:     "PNG image extraction",
			filename: "sample.png",
			wantErr:  false,
			expected: "A picture sample",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join("testdata", tt.filename)
			options := extractor.DefaultExtractOptions()

			extractorInstance := extractor.NewImageExtractor()
			result, err := extractorInstance.ExtractFromFile(filePath, options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageExtractor.ExtractFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Text != tt.expected {
					t.Errorf("ImageExtractor.ExtractFromFile() text = %v, expected %v", result.Text, tt.expected)
				}

				if result.FileType != "image" {
					t.Errorf("ImageExtractor.ExtractFromFile() fileType = %v, expected 'image'", result.FileType)
				}

				if result.Metadata == nil {
					t.Error("ImageExtractor.ExtractFromFile() metadata should not be nil")
				}

				t.Logf("Image extraction successful. Text: %s", result.Text)
				t.Logf("Metadata: %+v", result.Metadata)
			}
		})
	}
}

func TestImageExtractorSupportedTypes(t *testing.T) {
	extractorInstance := extractor.NewImageExtractor()

	supportedTypes := extractorInstance.SupportedTypes()
	expectedTypes := []string{"png", "jpg", "jpeg", "gif", "bmp", "tiff"}

	if len(supportedTypes) != len(expectedTypes) {
		t.Errorf("Expected %d supported types, got %d", len(expectedTypes), len(supportedTypes))
	}

	for _, expectedType := range expectedTypes {
		found := false
		for _, supportedType := range supportedTypes {
			if supportedType == expectedType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected type %s not found in supported types", expectedType)
		}
	}
}

func TestImageExtractorInvalidFormat(t *testing.T) {
	// Test with a non-image file (use sample.txt)
	filePath := filepath.Join("testdata", "sample.txt")
	options := extractor.DefaultExtractOptions()

	extractorInstance := extractor.NewImageExtractor()
	result, err := extractorInstance.ExtractFromFile(filePath, options)
	if err != nil {
		t.Errorf("ImageExtractor.ExtractFromFile() unexpected error = %v", err)
		return
	}

	// Should return empty text for non-image files with error in metadata
	if result.Text != "" {
		t.Errorf("ImageExtractor.ExtractFromFile() should return empty text for non-image files, got: %s", result.Text)
	}

	// Check that error is recorded in metadata
	if errorMsg, exists := result.Metadata["error"]; !exists || errorMsg != "invalid image format" {
		t.Errorf("ImageExtractor.ExtractFromFile() expected 'invalid image format' error in metadata")
	}

	t.Logf("Non-image file extraction result: %s", result.Text)
	t.Logf("Metadata: %+v", result.Metadata)
}
