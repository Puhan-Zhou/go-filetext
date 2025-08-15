package extractor

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

// XLSXExtractor handles Excel XLSX file text extraction
type XLSXExtractor struct {
	PlainTextExtractor
}

// NewXLSXExtractor creates a new XLSX extractor
func NewXLSXExtractor() *XLSXExtractor {
	return &XLSXExtractor{}
}

// Extract extracts text and metadata from an XLSX file
func (e *XLSXExtractor) Extract(reader io.Reader, options ExtractOptions) (*ExtractResult, error) {
	// Read all content first
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read XLSX content: %w", err)
	}

	// Get basic text extraction from content
	result, err := e.PlainTextExtractor.Extract(strings.NewReader(string(content)), options)
	if err != nil {
		return nil, err
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "xlsx_temp_*.xlsx")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write content to temp file
	if _, err := tempFile.Write(content); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tempFile.Close()

	// Parse XLSX file
	xlsxFile, err := xlsx.OpenFile(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to parse XLSX file: %w", err)
	}

	// Extract text from all sheets
	var extractedText strings.Builder
	sheetCount := 0
	totalRows := 0
	totalCells := 0

	for _, sheet := range xlsxFile.Sheets {
		sheetCount++

		// Iterate through rows using ForEachRow
		sheet.ForEachRow(func(row *xlsx.Row) error {
			totalRows++
			// Iterate through cells in the row
			row.ForEachCell(func(cell *xlsx.Cell) error {
				totalCells++
				cellText := strings.TrimSpace(cell.String())
				if cellText != "" {
					// Simply concatenate cell text without any separators
					extractedText.WriteString(cellText)
				}
				return nil
			})
			return nil
		})
	}

	// Update result with extracted text and XLSX-specific metadata
	result.Text = extractedText.String()
	result.FileType = "xlsx"
	result.Metadata["sheets"] = strconv.Itoa(sheetCount)
	result.Metadata["rows"] = strconv.Itoa(totalRows)
	result.Metadata["cells"] = strconv.Itoa(totalCells)
	result.Metadata["character_count"] = strconv.Itoa(len(result.Text))
	result.Metadata["line_count"] = strconv.Itoa(strings.Count(result.Text, "\n") + 1)

	return result, nil
}

// ExtractFromFile extracts text from an XLSX file
func (e *XLSXExtractor) ExtractFromFile(filePath string, options ExtractOptions) (*ExtractResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return e.Extract(file, options)
}
