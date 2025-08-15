package extractor

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// DOCXExtractor handles Microsoft Word DOCX file text extraction
type DOCXExtractor struct {
	PlainTextExtractor
}

// NewDOCXExtractor creates a new DOCX extractor
func NewDOCXExtractor() *DOCXExtractor {
	return &DOCXExtractor{}
}

// Extract extracts text and metadata from a DOCX file
func (e *DOCXExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	// Read all content first
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read DOCX content: %w", err)
	}

	// Get basic text extraction from content
	result, err := e.PlainTextExtractor.Extract(strings.NewReader(string(content)), options)
	if err != nil {
		return nil, err
	}

	// Extract text by manually parsing DOCX structure
	extractedText, err := e.extractTextFromDOCX(content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from DOCX file: %w", err)
	}

	// Update result with extracted text and DOCX-specific metadata
	result.Text = extractedText
	result.FileType = "docx"
	result.Metadata["paragraphs"] = strconv.Itoa(strings.Count(result.Text, "\n") + 1)
	result.Metadata["table_count"] = "0"
	result.Metadata["characters"] = strconv.Itoa(len(result.Text))
	result.Metadata["line_count"] = strconv.Itoa(strings.Count(result.Text, "\n") + 1)

	return result, nil
}

// ExtractFromFile extracts text from a DOCX file
func (e *DOCXExtractor) ExtractFromFile(filePath string, options ExtractOptions) (*ExtractResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return e.Extract(file, options)
}

// extractTextFromDOCX manually extracts text from DOCX content
func (e *DOCXExtractor) extractTextFromDOCX(content []byte) (string, error) {
	// Create a zip reader from the DOCX content
	reader := strings.NewReader(string(content))
	zipReader, err := zip.NewReader(reader, int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("failed to read DOCX as zip: %w", err)
	}

	// Find and read word/document.xml
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open document.xml: %w", err)
			}
			xmlContent, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return "", fmt.Errorf("failed to read document.xml: %w", err)
			}

			// Extract text from XML content
			return e.extractTextFromXML(string(xmlContent)), nil
		}
	}

	return "", fmt.Errorf("document.xml not found in DOCX file")
}

// extractTextFromXML extracts text content from Word document XML while preserving paragraph structure
func (e *DOCXExtractor) extractTextFromXML(xmlContent string) string {
	var textBuilder strings.Builder

	// Use regex to find paragraphs and extract text in order
	paragraphRegex := regexp.MustCompile(`<w:p[^>]*>(.*?)</w:p>`)
	textRegex := regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)

	// Find all paragraphs
	paragraphs := paragraphRegex.FindAllStringSubmatch(xmlContent, -1)

	for i, paragraph := range paragraphs {
		if len(paragraph) > 1 {
			// Extract all text runs from this paragraph
			textRuns := textRegex.FindAllStringSubmatch(paragraph[1], -1)
			
			for _, textRun := range textRuns {
				if len(textRun) > 1 {
					textBuilder.WriteString(textRun[1])
				}
			}
			
			// Add newline after each paragraph (except the last one)
			if i < len(paragraphs)-1 {
				textBuilder.WriteString("\n")
			}
		}
	}

	return strings.TrimSpace(textBuilder.String())
}