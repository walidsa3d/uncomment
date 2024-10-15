package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputFile := flag.String("input", "", "Input configuration file")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Please provide an input file using the -input flag")
		os.Exit(1)
	}

	err := processFile(*inputFile)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File processed successfully")
}

func processFile(filename string) error {
	// Create backup file
	backupFile := filename + ".bak"
	err := copyFile(filename, backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	// Open input file
	input, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer input.Close()

	// Create temporary output file
	tempFile := filename + ".temp"
	output, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer output.Close()

	// Process the file
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)

	for scanner.Scan() {
		line := scanner.Text()
		if !isComment(line) {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return fmt.Errorf("failed to write to temporary file: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}

	writer.Flush()

	// Replace original file with processed file
	err = os.Rename(tempFile, filename)
	if err != nil {
		return fmt.Errorf("failed to replace original file: %v", err)
	}

	return nil
}

func isComment(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return strings.HasPrefix(trimmedLine, "#") || strings.HasPrefix(trimmedLine, "//")
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}