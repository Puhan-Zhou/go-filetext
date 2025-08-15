package extractor

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// PPTXExtractor handles Microsoft PowerPoint PPTX file text extraction
type PPTXExtractor struct {
	PlainTextExtractor
}

// NewPPTXExtractor creates a new PPTX extractor
func NewPPTXExtractor() *PPTXExtractor {
	return &PPTXExtractor{}
}

// Extract extracts text and metadata from a PPTX file
func (e *PPTXExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	// Read all content first
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read PPTX content: %w", err)
	}

	// Extract text by manually parsing PPTX structure
	extractedText, slideCount, err := e.extractTextFromPPTX(content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from PPTX file: %w", err)
	}

	// Create result with extracted text and PPTX-specific metadata
	result := &ExtractResult{
		Text:     extractedText,
		FileType: "pptx",
		Metadata: make(map[string]interface{}),
	}

	result.Metadata["slides"] = strconv.Itoa(slideCount)
	result.Metadata["characters"] = strconv.Itoa(len(result.Text))
	result.Metadata["line_count"] = strconv.Itoa(strings.Count(result.Text, "\n") + 1)

	return result, nil
}

// ExtractFromFile extracts text from a PPTX file
func (e *PPTXExtractor) ExtractFromFile(filePath string, options ExtractOptions) (*ExtractResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return e.Extract(file, options)
}

// extractTextFromPPTX manually extracts text from PPTX content
func (e *PPTXExtractor) extractTextFromPPTX(content []byte) (string, int, error) {
	// Create a zip reader from the content
	zipReader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", 0, fmt.Errorf("failed to read PPTX as ZIP: %w", err)
	}

	var allText []string
	slideCount := 0

	// Look for slide files in the ZIP archive
	for _, file := range zipReader.File {
		// PPTX slides are stored as ppt/slides/slideX.xml
		if strings.HasPrefix(file.Name, "ppt/slides/slide") && strings.HasSuffix(file.Name, ".xml") {
			slideCount++
			
			// Open and read the slide XML file
			rc, err := file.Open()
			if err != nil {
				continue
			}

			xmlContent, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			// Extract text from the XML content
			slideText := e.extractTextFromSlideXML(string(xmlContent))
			if slideText != "" {
				allText = append(allText, slideText)
			}
		}
	}

	return strings.Join(allText, "\n\n"), slideCount, nil
}

// extractTextFromSlideXML extracts text content from slide XML
func (e *PPTXExtractor) extractTextFromSlideXML(xmlContent string) string {
	// PPTX text is stored in <a:t> tags within the slide XML
	re := regexp.MustCompile(`<a:t[^>]*>([^<]*)</a:t>`)
	matches := re.FindAllStringSubmatch(xmlContent, -1)

	var textParts []string
	for _, match := range matches {
		if len(match) > 1 && strings.TrimSpace(match[1]) != "" {
			textParts = append(textParts, strings.TrimSpace(match[1]))
		}
	}

	return strings.Join(textParts, " ")
}