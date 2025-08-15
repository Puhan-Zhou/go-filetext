package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Puhan-Zhou/go-file-plain-text/pkg/extractor"
)

func main() {
	// Check if filename is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		fmt.Println("Example: go run main.go document.pdf")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Create factory and extractor based on file extension
	factory := extractor.NewExtractorFactory()
	extractorInstance, err := factory.CreateExtractorFromPath(filename)
	if err != nil {
		log.Fatalf("Error creating extractor for %s: %v", filename, err)
	}

	// Extract text using default options
	options := extractor.DefaultExtractOptions()
	result, err := extractorInstance.ExtractFromFile(filename, options)
	if err != nil {
		log.Fatalf("Error extracting text from %s: %v", filename, err)
	}

	// Print the extracted text
	fmt.Print(result.Text)
}
