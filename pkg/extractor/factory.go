package extractor

import (
	"fmt"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// ExtractorFactory creates appropriate extractors based on file type
type ExtractorFactory struct {
}

// NewExtractorFactory creates a new extractor factory
func NewExtractorFactory() *ExtractorFactory {
	return &ExtractorFactory{}
}

// CreateExtractorFromPath creates an extractor based on file path using filetype detection
func (f *ExtractorFactory) CreateExtractorFromPath(filePath string) (TextExtractor, error) {
	// Use filetype module for content-based detection
	mtype, err := mimetype.DetectFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to detect file type: %w", err)
	}

	// If filetype detected something, use it
	if mtype.String() != "unknown" {
		// Map MIME types and extensions to our extractors
		mimeType := mtype.String()

		// Handle specific MIME types first
		switch mimeType {
		case "application/pdf":
			return NewPDFExtractor(), nil
		case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
			return NewDOCXExtractor(), nil
		case "application/msword":
			return NewLegacyDOCExtractor(), nil // Legacy DOC format
		case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
			return NewXLSXExtractor(), nil
		case "application/vnd.ms-excel":
			return NewLegacyXLSExtractor(), nil // Legacy XLS format
		case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
			return NewPPTXExtractor(), nil
		case "application/vnd.ms-powerpoint":
			return NewLegacyPPTExtractor(), nil // Legacy PPT format
		default:
			if strings.HasPrefix(mimeType, "text/") {
				return NewPlainTextExtractor(), nil
			} else {
				return nil, fmt.Errorf("unsupported file type %s for %s", mimeType, filePath)
			}
		}
	} else {
		return nil, fmt.Errorf("unknown file type for: %s", filePath)
	}
}
