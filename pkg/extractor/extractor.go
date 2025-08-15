package extractor

import (
	"io"
	"time"
)

// TextExtractor defines the interface for extracting text from various file formats
type TextExtractor interface {
	// Extract extracts text from an io.Reader
	Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error)
	
	// ExtractFromFile extracts text from a file path
	ExtractFromFile(filepath string, options ExtractOptions) (*ExtractResult, error)
	
	// SupportedTypes returns a list of supported file types/extensions
	SupportedTypes() []string
}

// ExtractOptions contains configuration options for text extraction
type ExtractOptions struct {
	// FileType overrides automatic file type detection
	FileType string
	
	// OCRLanguage specifies the language for OCR processing (for images)
	OCRLanguage string
	
	// MaxFileSize sets the maximum file size to process (in bytes)
	MaxFileSize int64
	
	// Timeout sets the maximum time to spend on extraction
	Timeout time.Duration
	
	// PreserveFormatting indicates whether to preserve text formatting
	PreserveFormatting bool
}

// ExtractResult contains the result of text extraction
type ExtractResult struct {
	// Text is the extracted plain text content
	Text string
	
	// Metadata contains additional information about the extraction
	Metadata map[string]interface{}
	
	// FileType is the detected or specified file type
	FileType string
	
	// ProcessingTime is the time taken for extraction
	ProcessingTime time.Duration
}

// ExtractorError represents errors that occur during text extraction
type ExtractorError struct {
	Message   string
	FileType  string
	Operation string
	Cause     error
}

func (e *ExtractorError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *ExtractorError) Unwrap() error {
	return e.Cause
}

// NewExtractorError creates a new ExtractorError
func NewExtractorError(message, fileType, operation string, cause error) *ExtractorError {
	return &ExtractorError{
		Message:   message,
		FileType:  fileType,
		Operation: operation,
		Cause:     cause,
	}
}

// DefaultExtractOptions returns default extraction options
func DefaultExtractOptions() ExtractOptions {
	return ExtractOptions{
		MaxFileSize:        100 * 1024 * 1024, // 100MB
		Timeout:           30 * time.Second,
		OCRLanguage:       "eng",
		PreserveFormatting: false,
	}
}