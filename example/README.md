# Simple Text Extractor

A simple command line program that extracts text from various file formats and prints it to stdout.

## Usage

```bash
go run main.go <filename>
```

## Examples

```bash
# Extract text from a PDF file
go run main.go document.pdf

# Extract text from a Word document
go run main.go report.docx

# Extract text from a PowerPoint presentation
go run main.go presentation.pptx

# Extract text from an Excel spreadsheet
go run main.go data.xlsx

# Extract text from a plain text file
go run main.go notes.txt

# Extract text from a CSV file
go run main.go data.csv

# Extract text from a Markdown file
go run main.go README.md
```

## Supported File Types

- **Text files**: .txt, .text
- **CSV files**: .csv
- **Markdown files**: .md, .markdown
- **PDF files**: .pdf
- **Microsoft Word**: .docx
- **Microsoft Excel**: .xlsx
- **Microsoft PowerPoint**: .pptx, .ppt

## Building

To build a standalone executable:

```bash
go build -o text-extractor main.go
./text-extractor document.pdf
```

## Features

- Automatic file type detection based on file extension
- Clean text output suitable for piping to other commands
- Error handling with helpful messages
- Support for multiple document formats

## Error Handling

The program will exit with status code 1 and display an error message if:
- No filename is provided
- The file cannot be opened
- The file format is not supported
- Text extraction fails