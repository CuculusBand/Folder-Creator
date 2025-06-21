package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

// FileProcessor is a struct that holds the file processing logic
type FileProcessor struct {
	TableFilePath string
	DestPath      string
	TableData     [][]string
}

// Create new FileProcessor instance
func NewFileProcessor() *FileProcessor {
	return &FileProcessor{}
}

// load CSV file
func (p *FileProcessor) ReadCSVFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	// Allow variable number of fields
	reader.FieldsPerRecord = -1
	return reader.ReadAll()
}

// load XLSX file
func (p *FileProcessor) ReadXLSXFile(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Read the first sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("did not find any sheets in the file")
	}
	// Read all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	// Ensure all rows have the same number of columns
	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	for i := range rows {
		for len(rows[i]) < maxCols {
			rows[i] = append(rows[i], "")
		}
	}
	return rows, nil
}

// Load slected file
func (p *FileProcessor) LoadFile(filePath string) error {
	p.TableFilePath = filePath
	ext := strings.ToLower(filepath.Ext(filePath))

	var data [][]string
	var err error

	switch ext {
	case ".csv":
		data, err = p.ReadCSVFile(filePath)
	case ".xlsx":
		data, err = p.ReadXLSXFile(filePath)
	default:
		err = fmt.Errorf("file not supported: %s", ext)
	}
	if err != nil {
		return err
	}
	p.TableData = data
	return nil
}

// Create folders based on the loaded table
func (p *FileProcessor) GenerateFolders() (int, error) {
	successCount := 0
	for _, row := range p.TableData {
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}

		// Create the first level folder
		level1Path := filepath.Join(p.DestPath, strings.TrimSpace(row[0]))
		if err := os.MkdirAll(level1Path, 0755); err != nil {
			return successCount, fmt.Errorf("falied to create %s: %v", row[0], err)
		}
		successCount++

		// Create subfolders if they exist
		for i := 1; i < len(row); i++ {
			if strings.TrimSpace(row[i]) == "" {
				continue
			}

			subPath := filepath.Join(level1Path, strings.TrimSpace(row[i]))
			if err := os.MkdirAll(subPath, 0755); err != nil {
				return successCount, fmt.Errorf("falied to create %s: %v", row[i], err)
			}
			successCount++
		}
	}

	return successCount, nil
}

// Clear all content in the processor
func (p *FileProcessor) Clear() {
	p.TableFilePath = ""
	p.DestPath = ""
	p.TableData = [][]string{}
}
