// This file is part of csv-analysis.
//
// Copyright (C) 2017  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package csvutil provides ways to extract csv data from one or more files.
*/
package csvutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// CSVFiles - Struct containing CSV file information.
type CSVFiles struct {
	Files []string
	// Indicates if the CSV files don't have a header.
	NoHeader bool
	// Indicates if the 0 value should be filtered.
	FilterZero bool
}

// New - Returns a `*csvutil.CSVFiles` with the files given.
func New(files ...string) *CSVFiles {
	return &CSVFiles{Files: files, NoHeader: false, FilterZero: false}
}

// GetCSVColumns - Reads csv lines from *csvutil.CSVFiles and returns the requested columns.
// If the column lenghts are different, it will error out (necessary here? maybe the caller should check for that.)
func (cf *CSVFiles) GetCSVColumns(columns ...int) ([][]string, error) {
	columnsData := make([][]string, len(columns))
	for _, file := range cf.Files {
		fh, err := os.Open(file)
		if err != nil {
			return columnsData, err
		}
		defer fh.Close()
		fs, err := getCSVColumns(fh, columns...)
		fh.Close()
		if err != nil {
			return columnsData, err
		}
		l := len(fs[0])
		for i, columnString := range fs {
			lc := len(columnString)
			if l == 0 {
				fmt.Fprintf(os.Stderr, "Column %d is empty, file: %s\n", columns[i], file)
				continue
			}
			if l != lc {
				return nil, fmt.Errorf("Column lenghts do not match")
			}
			if cf.NoHeader {
				columnsData[i] = append(columnsData[i], columnString...)
			} else {
				columnsData[i] = append(columnsData[i], columnString[1:]...)
			}
		}
	}
	return columnsData, nil
}

// PrintCSVRows - prints the given csv rows
func (cf *CSVFiles) PrintCSVRows(rows ...int) error {
	for _, file := range cf.Files {
		fh, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fh.Close()
		rowData, err := getCSVRows(fh, rows...)
		fh.Close()
		if err != nil {
			return err
		}
		for i, row := range rowData {
			if rows[i] == 1 {
				fmt.Printf("Header Row\n")
			} else {
				fmt.Printf("Row %d\n", rows[i])
			}
			for i, e := range row {
				fmt.Printf("%d: %s\n", i+1, e)
			}
		}
	}
	return nil
}

// GetFloat64Columns - given a set of CSV files and a list of columns, it will return those columns as a slice of floats.
// If filterZero is set, it will ignore Zero values.
// If the column lenghts are different, it will error out (necessary here? maybe the caller should check for that.)
func (cf *CSVFiles) GetFloat64Columns(columns ...int) ([][]float64, error) {
	cSlices, err := cf.GetCSVColumns(columns...)
	if err != nil {
		return nil, err
	}
	sliceDatasets := make([][]float64, len(columns))
	for i, cSlice := range cSlices {
		for j := range cSlice {
			trimmed := strings.TrimSpace(cSlice[j])
			x64, err := strconv.ParseFloat(trimmed, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			} else {
				if !cf.FilterZero || (cf.FilterZero && x64 != 0) {
					sliceDatasets[i] = append(sliceDatasets[i], x64)
				}
			}
		}
	}
	return sliceDatasets, nil
}

// getCSVRows - Reads csv lines from `reader` and returns the requested rows.
func getCSVRows(reader io.Reader, rows ...int) ([][]string, error) {
	rowsData := make([][]string, len(rows))
	// Verify query
	maxRow := 0
	for _, r := range rows {
		if r <= 0 {
			return nil, fmt.Errorf("Row index error: %d <= 0!", r)
		}
		if maxRow < r {
			maxRow = r
		}
	}
	r := csv.NewReader(reader)
	r.FieldsPerRecord = -1
	rowCounter := 0
	for {
		rowCounter++
		if rowCounter > maxRow {
			break
		}
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return nil, err
		}
		for i, row := range rows {
			if rowCounter == row {
				rowsData[i] = append(rowsData[i], record...)
			}
		}
	}
	return rowsData, nil
}

// getCSVColumns - Reads csv lines from `reader` and returns the requested columns.
func getCSVColumns(reader io.Reader, columns ...int) ([][]string, error) {
	columnsData := make([][]string, len(columns))
	// Verify query
	for _, c := range columns {
		if c <= 0 {
			return nil, fmt.Errorf("Column index error: %d <= 0!", c)
		}
	}
	r := csv.NewReader(reader)
	r.FieldsPerRecord = -1
	row := 0
	for {
		row++
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			return nil, err
		}
		for i, column := range columns {
			if len(record) >= column {
				columnsData[i] = append(columnsData[i], record[column-1])
			}
		}
	}
	return columnsData, nil
}
