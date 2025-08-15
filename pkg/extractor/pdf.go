package extractor

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dslipak/pdf"
)

// PDFExtractor provides PDF text extraction capabilities
type PDFExtractor struct{}

// NewPDFExtractor creates a new PDF extractor
func NewPDFExtractor() *PDFExtractor {
	return &PDFExtractor{}
}

// Extract extracts text from a PDF reader
func (e *PDFExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	start := time.Now()

	// Read all content into memory
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, NewExtractorError("failed to read PDF content", "pdf", "read", err)
	}

	// Check file size limit
	if options.MaxFileSize > 0 && int64(len(content)) > options.MaxFileSize {
		return nil, NewExtractorError(
			fmt.Sprintf("file size %d exceeds limit %d", len(content), options.MaxFileSize),
			"pdf", "size_check", nil)
	}

	// Create a temporary file for the PDF library
	tempFile, err := os.CreateTemp("", "pdf_extract_*.pdf")
	if err != nil {
		return nil, NewExtractorError("failed to create temp file", "pdf", "temp_file", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write content to temp file
	if _, err := tempFile.Write(content); err != nil {
		return nil, NewExtractorError("failed to write temp file", "pdf", "temp_write", err)
	}
	tempFile.Close()

	// Open PDF file
	r, err := pdf.Open(tempFile.Name())
	if err != nil {
		return nil, NewExtractorError("failed to open PDF", "pdf", "open", err)
	}

	// Extract text from all pages
	var textBuilder strings.Builder
	pageCount := r.NumPage()

	for pageNum := 1; pageNum <= pageCount; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		// Extract text from page
		pageText, err := e.extractPageText(page)
		if err != nil {
			// Log error but continue with other pages
			continue
		}

		textBuilder.WriteString(pageText)
		if pageNum < pageCount {
			textBuilder.WriteString("\n\n") // Separate pages
		}
	}

	text := textBuilder.String()

	// Normalize line endings if not preserving formatting
	if !options.PreserveFormatting {
		text = e.normalizeLineEndings(text)
	}

	metadata := map[string]interface{}{
		"page_count": pageCount,
		"size_bytes": len(content),
		"line_count": strings.Count(text, "\n") + 1,
		"char_count": len([]rune(text)),
		"word_count": len(strings.Fields(text)),
	}

	return &ExtractResult{
		Text:           text,
		Metadata:       metadata,
		FileType:       "pdf",
		ProcessingTime: time.Since(start),
	}, nil
}

// ExtractFromFile extracts text from a PDF file
func (e *PDFExtractor) ExtractFromFile(filepath string, options ExtractOptions) (*ExtractResult, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, NewExtractorError("failed to open file", "pdf", "open", err)
	}
	defer file.Close()

	return e.Extract(file, options)
}

// SupportedTypes returns the file types supported by this extractor
func (e *PDFExtractor) SupportedTypes() []string {
	return []string{"pdf"}
}

// extractPageText extracts text from a single PDF page
func (e *PDFExtractor) extractPageText(page pdf.Page) (string, error) {
	// Get page content
	content := page.Content()
	if content.Text == nil {
		return "", nil
	}

	var textBuilder strings.Builder
	for _, text := range content.Text {
		// Simply concatenate all text elements without adding extra spaces
		// The PDF library already provides proper spacing as separate elements
		textBuilder.WriteString(text.S)
	}

	return strings.TrimSpace(textBuilder.String()), nil
}

// normalizeLineEndings converts different line ending formats to \n
func (e *PDFExtractor) normalizeLineEndings(text string) string {
	// Replace Windows line endings
	text = strings.ReplaceAll(text, "\r\n", "\n")
	// Replace Mac line endings
	text = strings.ReplaceAll(text, "\r", "\n")
	return text
}
