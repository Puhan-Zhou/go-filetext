package extractor

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// LegacyDOCExtractor handles legacy Microsoft Word DOC files
// Note: This is a placeholder that provides informative error messages
// Legacy DOC files use a proprietary binary format that requires specialized libraries
type LegacyDOCExtractor struct {
	PlainTextExtractor
}

// NewLegacyDOCExtractor creates a new legacy DOC extractor
func NewLegacyDOCExtractor() *LegacyDOCExtractor {
	return &LegacyDOCExtractor{}
}

// Extract extracts text from a legacy DOC file
func (e *LegacyDOCExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	return nil, fmt.Errorf("legacy DOC format is not supported: DOC files use a proprietary binary format that requires specialized libraries. Please convert to DOCX format for text extraction")
}

// ExtractFromFile extracts text from a legacy DOC file
func (e *LegacyDOCExtractor) ExtractFromFile(filePath string, options ExtractOptions) (*ExtractResult, error) {
	return nil, fmt.Errorf("legacy DOC format is not supported: %s uses a proprietary binary format. Please convert to DOCX format for text extraction", filepath.Base(filePath))
}

// SupportedTypes returns the file types supported by this extractor
func (e *LegacyDOCExtractor) SupportedTypes() []string {
	return []string{"doc"}
}

// LegacyXLSExtractor handles legacy Microsoft Excel XLS files
// Note: This is a placeholder that provides informative error messages
// Legacy XLS files use a proprietary binary format that requires specialized libraries
type LegacyXLSExtractor struct {
	PlainTextExtractor
}

// NewLegacyXLSExtractor creates a new legacy XLS extractor
func NewLegacyXLSExtractor() *LegacyXLSExtractor {
	return &LegacyXLSExtractor{}
}

// Extract extracts text from a legacy XLS file
func (e *LegacyXLSExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	return nil, fmt.Errorf("legacy XLS format is not supported: XLS files use a proprietary binary format that requires specialized libraries. Please convert to XLSX format for text extraction")
}

// ExtractFromFile extracts text from a legacy XLS file
func (e *LegacyXLSExtractor) ExtractFromFile(filePath string, options ExtractOptions) (*ExtractResult, error) {
	return nil, fmt.Errorf("legacy XLS format is not supported: %s uses a proprietary binary format. Please convert to XLSX format for text extraction", filepath.Base(filePath))
}

// SupportedTypes returns the file types supported by this extractor
func (e *LegacyXLSExtractor) SupportedTypes() []string {
	return []string{"xls"}
}

// LegacyPPTExtractor handles legacy Microsoft PowerPoint PPT files
// Note: This is a placeholder that provides informative error messages
// Legacy PPT files use a proprietary binary format that requires specialized libraries
type LegacyPPTExtractor struct {
	PlainTextExtractor
}

// NewLegacyPPTExtractor creates a new legacy PPT extractor
func NewLegacyPPTExtractor() *LegacyPPTExtractor {
	return &LegacyPPTExtractor{}
}

// Extract extracts text from a legacy PPT file
func (e *LegacyPPTExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	return nil, fmt.Errorf("legacy PPT format is not supported: PPT files use a proprietary binary format that requires specialized libraries. Please convert to PPTX format for text extraction")
}

// ExtractFromFile extracts text from a legacy PPT file
func (e *LegacyPPTExtractor) ExtractFromFile(filePath string, options ExtractOptions) (*ExtractResult, error) {
	return nil, fmt.Errorf("legacy PPT format is not supported: %s uses a proprietary binary format. Please convert to PPTX format for text extraction", filepath.Base(filePath))
}

// SupportedTypes returns the file types supported by this extractor
func (e *LegacyPPTExtractor) SupportedTypes() []string {
	return []string{"ppt"}
}

// Helper function to detect if a file is a legacy Office format
func IsLegacyOfficeFormat(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".doc" || ext == ".xls" || ext == ".ppt"
}

// Helper function to suggest modern equivalent
func GetModernEquivalent(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".doc":
		return "DOCX"
	case ".xls":
		return "XLSX"
	case ".ppt":
		return "PPTX"
	default:
		return "unknown"
	}
}