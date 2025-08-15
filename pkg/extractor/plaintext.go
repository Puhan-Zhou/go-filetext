package extractor

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// PlainTextExtractor handles extraction from plain text files
type PlainTextExtractor struct{}

// NewPlainTextExtractor creates a new plain text extractor
func NewPlainTextExtractor() *PlainTextExtractor {
	return &PlainTextExtractor{}
}

// Extract extracts text from a plain text reader
func (e *PlainTextExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	start := time.Now()

	// Read all content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, NewExtractorError("failed to read content", "plaintext", "read", err)
	}

	// Check file size limit
	if options.MaxFileSize > 0 && int64(len(content)) > options.MaxFileSize {
		return nil, NewExtractorError(
			fmt.Sprintf("file size %d exceeds limit %d", len(content), options.MaxFileSize),
			"plaintext", "size_check", nil)
	}

	// Detect and convert encoding
	text, encoding, err := e.detectAndConvertEncoding(content)
	if err != nil {
		return nil, NewExtractorError("failed to convert encoding", "plaintext", "encoding", err)
	}

	// Normalize line endings if not preserving formatting
	if !options.PreserveFormatting {
		text = e.normalizeLineEndings(text)
	}

	metadata := map[string]interface{}{
		"encoding":   encoding,
		"size_bytes": len(content),
		"line_count": strings.Count(text, "\n") + 1,
		"char_count": len([]rune(text)),
	}

	return &ExtractResult{
		Text:           text,
		Metadata:       metadata,
		FileType:       "plaintext",
		ProcessingTime: time.Since(start),
	}, nil
}

// ExtractFromFile extracts text from a plain text file
func (e *PlainTextExtractor) ExtractFromFile(filepath string, options ExtractOptions) (*ExtractResult, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, NewExtractorError("failed to open file", "plaintext", "open", err)
	}
	defer file.Close()

	return e.Extract(file, options)
}

// SupportedTypes returns the file types supported by this extractor
func (e *PlainTextExtractor) SupportedTypes() []string {
	return []string{"txt", "csv", "yaml", "yml", "json", "xml", "md", "markdown", "log", "conf", "cfg", "ini"}
}

// detectAndConvertEncoding detects the encoding and converts to UTF-8
func (e *PlainTextExtractor) detectAndConvertEncoding(content []byte) (string, string, error) {
	// First check if it's already valid UTF-8
	if utf8.Valid(content) {
		return string(content), "UTF-8", nil
	}

	// Try common encodings
	encodings := []struct {
		name string
		enc  encoding.Encoding
	}{
		{"UTF-16LE", unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)},
		{"UTF-16BE", unicode.UTF16(unicode.BigEndian, unicode.UseBOM)},
		{"Windows-1252", charmap.Windows1252},
		{"ISO-8859-1", charmap.ISO8859_1},
		{"Windows-1251", charmap.Windows1251},
	}

	for _, enc := range encodings {
		decoder := enc.enc.NewDecoder()
		result, err := io.ReadAll(transform.NewReader(strings.NewReader(string(content)), decoder))
		if err == nil && utf8.Valid(result) {
			return string(result), enc.name, nil
		}
	}

	// If all else fails, replace invalid UTF-8 sequences
	text := strings.ToValidUTF8(string(content), "�")
	return text, "UTF-8 (with replacements)", nil
}

// normalizeLineEndings converts all line endings to \n
func (e *PlainTextExtractor) normalizeLineEndings(text string) string {
	// Replace Windows (\r\n) and Mac (\r) line endings with Unix (\n)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return text
}

// CSVExtractor provides specialized CSV handling
type CSVExtractor struct {
	*PlainTextExtractor
}

// NewCSVExtractor creates a new CSV extractor
func NewCSVExtractor() *CSVExtractor {
	return &CSVExtractor{
		PlainTextExtractor: NewPlainTextExtractor(),
	}
}

// Extract extracts text from CSV with additional metadata
func (e *CSVExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	result, err := e.PlainTextExtractor.Extract(reader, options)
	if err != nil {
		return nil, err
	}

	// Add CSV-specific metadata
	lines := strings.Split(result.Text, "\n")
	if len(lines) > 0 {
		firstLine := lines[0]
		columnCount := strings.Count(firstLine, ",") + 1
		result.Metadata["column_count"] = columnCount
		result.Metadata["has_header"] = e.detectCSVHeader(lines)
	}

	// Override the file type set by parent extractor
	result.FileType = "csv"
	return result, nil
}

// ExtractFromFile extracts text from a CSV file
func (e *CSVExtractor) ExtractFromFile(filepath string, options ExtractOptions) (*ExtractResult, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, NewExtractorError("failed to open file", "csv", "open", err)
	}
	defer file.Close()

	return e.Extract(file, options)
}

// detectCSVHeader tries to determine if the first row is a header
func (e *CSVExtractor) detectCSVHeader(lines []string) bool {
	if len(lines) < 2 {
		return false
	}

	// Simple heuristic: if first row has fewer numbers than second row, it's likely a header
	firstRowNumbers := e.countNumbers(lines[0])
	secondRowNumbers := e.countNumbers(lines[1])

	return firstRowNumbers < secondRowNumbers
}

// countNumbers counts numeric values in a CSV row
func (e *CSVExtractor) countNumbers(row string) int {
	fields := strings.Split(row, ",")
	numberCount := 0

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field != "" {
			// Simple check for numeric content
			if strings.ContainsAny(field, "0123456789") {
				numberCount++
			}
		}
	}

	return numberCount
}

// SupportedTypes returns CSV-specific supported types
func (e *CSVExtractor) SupportedTypes() []string {
	return []string{"csv", "tsv"}
}

// MarkdownExtractor provides specialized Markdown handling
type MarkdownExtractor struct {
	*PlainTextExtractor
}

// NewMarkdownExtractor creates a new Markdown extractor
func NewMarkdownExtractor() *MarkdownExtractor {
	return &MarkdownExtractor{
		PlainTextExtractor: NewPlainTextExtractor(),
	}
}

// Extract extracts text from Markdown with optional formatting removal
func (e *MarkdownExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	result, err := e.PlainTextExtractor.Extract(reader, options)
	if err != nil {
		return nil, err
	}

	// If not preserving formatting, strip Markdown syntax
	if !options.PreserveFormatting {
		result.Text = e.stripMarkdownSyntax(result.Text)
	}

	// Add Markdown-specific metadata
	result.Metadata["heading_count"] = e.countHeadings(result.Text)
	result.Metadata["link_count"] = e.countLinks(result.Text)

	// Override the file type set by parent extractor
	result.FileType = "markdown"

	return result, nil
}

// ExtractFromFile extracts text from a Markdown file
func (e *MarkdownExtractor) ExtractFromFile(filepath string, options ExtractOptions) (*ExtractResult, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, NewExtractorError("failed to open file", "markdown", "open", err)
	}
	defer file.Close()

	return e.Extract(file, options)
}

// stripMarkdownSyntax removes common Markdown formatting
func (e *MarkdownExtractor) stripMarkdownSyntax(text string) string {
	lines := strings.Split(text, "\n")
	var cleanLines []string

	for _, line := range lines {
		// Remove heading markers
		line = strings.TrimLeft(line, "# ")

		// Remove bold and italic markers
		line = strings.ReplaceAll(line, "**", "")
		line = strings.ReplaceAll(line, "__", "")
		line = strings.ReplaceAll(line, "*", "")
		line = strings.ReplaceAll(line, "_", "")

		// Remove code markers
		line = strings.ReplaceAll(line, "`", "")

		// Remove list markers (simple approach)
		line = strings.TrimLeft(line, "- + ")

		cleanLines = append(cleanLines, line)
	}

	return strings.Join(cleanLines, "\n")
}

// countHeadings counts the number of headings in Markdown
func (e *MarkdownExtractor) countHeadings(text string) int {
	lines := strings.Split(text, "\n")
	count := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			count++
		}
	}

	return count
}

// countLinks counts the number of links in Markdown
func (e *MarkdownExtractor) countLinks(text string) int {
	// Simple count of [text](url) patterns
	count := 0
	inLink := false
	inURL := false

	for i, char := range text {
		switch char {
		case '[':
			if !inLink && !inURL {
				inLink = true
			}
		case ']':
			if inLink && i+1 < len(text) && text[i+1] == '(' {
				inLink = false
				inURL = true
			}
		case ')':
			if inURL {
				inURL = false
				count++
			}
		}
	}

	return count
}

// SupportedTypes returns Markdown-specific supported types
func (e *MarkdownExtractor) SupportedTypes() []string {
	return []string{"md", "markdown"}
}
