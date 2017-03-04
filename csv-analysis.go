// This file is part of csv-analysis.
//
// Copyright (C) 2017  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package main provides ways to analyse csv data from one or more files and generate basic statistic analysis.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/davidgamba/csv-analysis/csvutil"
	"github.com/davidgamba/csv-analysis/regression"
	"github.com/davidgamba/csv-analysis/stat"
	"github.com/davidgamba/go-getoptions"
	"github.com/gonum/matrix/mat64"
	// "github.com/gonum/optimize"
)

// Usage options

// noHeader - The csv file to be read has no header.
var noHeader bool

// filterZero - Ignore 0 value entries from statistical calculations.
var filterZero bool

// printError - prints the given error to STDERR.
func printError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
}

// printCSVColumnStats - Given a column and a set of csv files, it will print the statistical information for that column.
func printCSVColumnStats(files []string, column int) error {
	var fieldSliceDataset []float64

	for _, file := range files {
		cf := csvutil.New(file)
		cf.NoHeader = noHeader
		cf.FilterZero = filterZero
		fs, err := cf.GetFloat64Columns(column)
		if err != nil {
			return err
		}
		l := len(fs[0])
		if l == 0 {
			continue
		}
		fmt.Printf("Data: %d columns, %v\n", len(fs[0]), fs[0])
		fieldSliceDataset = append(fieldSliceDataset, fs[0]...)
	}
	stat.PrintSliceStats(fieldSliceDataset)
	return nil
}

func validateMinInt(min, value int) error {
	if value < min {
		return fmt.Errorf("can not be less than %d", min)
	}
	return nil
}
func validateMaxInt(max, value int) error {
	if value > max {
		return fmt.Errorf("can not be bigger than %d", max)
	}
	return nil
}

// trimSlice - It will trim the start and end of a slice.
func trimSlice(x []float64, trimStart, trimEnd int) ([]float64, error) {
	n := len(x)
	err := validateMinInt(0, trimStart)
	if err != nil {
		return nil, fmt.Errorf("trimStart %s", err)
	}
	err = validateMaxInt(n-1, trimStart)
	if err != nil {
		return nil, fmt.Errorf("trimStart %s", err)
	}
	err = validateMinInt(0, trimEnd)
	if err != nil {
		return nil, fmt.Errorf("trimEnd %s", err)
	}
	err = validateMaxInt(n-1, trimEnd)
	if err != nil {
		return nil, fmt.Errorf("trimEnd %s", err)
	}
	err = validateMaxInt(n-1, trimStart+trimEnd)
	if err != nil {
		return nil, fmt.Errorf("trimStart plus trimEnd %s", err)
	}
	return x[trimStart : n-trimEnd], nil
}

func synopsis() {
	synopsis := `csv-analysis --column|-c <n> <csv-file>...
       [--no-header|--nh] [--filter-zero|--fz] 

# Regression analysis
csv-analysis -x <n> -y <n> <csv-file>...
       [--no-header|--nh] [--filter-zero|--fz]
			 [--trim-start|--ts <n>] [--trim-end|--te <n>]
			 [--degree] [--regression] [--review]
			 [--plot-title <title>] [--plot-x-label <label>] [--plot-y-label <label>]

# Time plot
csv-analysis -x <n> -y <n> <csv-file>... -xtime <timeformat>
       [--no-header|--nh] [--filter-zero|--fz]
			 [--trim-start|--ts <n>] [--trim-end|--te <n>]
			 [--plot-title <title>] [--plot-x-label <label>] [--plot-y-label <label>]

# Inspect data and exit
csv-analysis [--show-header|-s] [--show-data|--sd] <csv-file>...

csv-analysis [--help]

# --column: Column to use for statistical analysis. n starts at 1.
#
# --no-header: The csv file has no header.
#              It is assumed that it does by default.
#
# --filter-zero: Ignore zeroes from statistical analysis.
#
# --x, --y: columns to use for X and Y when doing regression analysis.
#
# --trim-start, --trim-end: Trim fields from the CSV dataset.
#
# --degree: polynomial regression degree.
#
# --review: Show linear transformation graphs.
#
# --show-header: Show the header of the first csv file and exit.
#
# --show-data: Show the header and the first row of the first csv file and exit.
#
# --debug: Show debug output.
#
# --xtime: Mon Jan 2 15:04:05.000 MST 2006
#          Examples:
#          --xtime '2006/01/02 15:04:05.000'
#          --xtime '2006-01-02 15:04:05.000'
`
	fmt.Fprintln(os.Stderr, synopsis)
}

func main() {
	var column, xColumn int // field to analize
	var trimStart, trimEnd, degree int
	var pTitle, pYLabel, pXLabel string
	var xTimeFormat string
	var review bool

	opt := getoptions.New()
	// General options
	opt.Bool("help", false)
	opt.Bool("debug", false)
	// CSV review options
	opt.Bool("show-data", false, "sd")
	opt.Bool("show-header", false, "s")
	// CSV parsing options
	opt.BoolVar(&noHeader, "no-header", false, "nh")
	opt.BoolVar(&filterZero, "filter-zero", false, "fz")
	opt.BoolVar(&review, "review", false)
	// CSV data indicators
	opt.IntVar(&column, "column", 1, "c")
	opt.IntVar(&xColumn, "x", 1)
	opt.StringVarOptional(&xTimeFormat, "xtime", time.RFC3339)
	yColumns := opt.IntSliceMulti("y", 1, 99)
	// CSV data trimming
	opt.IntVar(&trimStart, "trim-start", 0, "ts")
	opt.IntVar(&trimEnd, "trim-end", 0, "te")
	// Linear Regression degree
	opt.IntVar(&degree, "degree", 1, "degree")
	// Action
	opt.Bool("regression", false, "r")
	// Plot options
	opt.StringVar(&pTitle, "plot-title", "Data", "pt")
	opt.StringVar(&pXLabel, "plot-x-label", "", "px")
	opt.StringVar(&pYLabel, "plot-y-label", "", "py")
	remaining, err := opt.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if opt.Called("help") {
		synopsis()
		os.Exit(1)
	}
	if opt.Called("debug") {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
	log.Println(remaining)
	if len(remaining) < 1 {
		fmt.Fprintf(os.Stderr, "ERROR: Missing file\n")
		os.Exit(1)
	}

	// Inspect data and quit
	if opt.Called("show-header") || opt.Called("show-data") {
		var err error
		cf := csvutil.New(remaining[0])
		cf.NoHeader = noHeader
		cf.FilterZero = filterZero
		if opt.Called("show-data") {
			err = cf.PrintCSVRows(1, 2)
		} else {
			err = cf.PrintCSVRows(1)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		os.Exit(1)
	}
	if opt.Called("x") && opt.Called("y") && opt.Called("xtime") {
		cf := csvutil.New(remaining...)
		cf.NoHeader = noHeader
		cf.FilterZero = filterZero
		sliceDatasetsString, err := cf.GetCSVColumns(xColumn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		var xSliceDataset []float64
		for _, e := range sliceDatasetsString[0] {
			trimmed := strings.TrimSpace(e)
			trimmedXTimeFormat := strings.TrimSpace(xTimeFormat)
			t, err := time.Parse(trimmedXTimeFormat, trimmed)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: time format '%s': %s\n", xTimeFormat, err)
				continue
			}
			xSliceDataset = append(xSliceDataset, float64(t.Unix()))
		}
		xTrimmed, err := trimSlice(xSliceDataset, trimStart, trimEnd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}

		sliceDatasets, err := cf.GetFloat64Columns(*yColumns...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		var sYTrimmed [][]float64
		for _, ySliceDataset := range sliceDatasets {
			yTrimmed, _ := trimSlice(ySliceDataset, trimStart, trimEnd)
			sYTrimmed = append(sYTrimmed, yTrimmed)
		}

		// TODO: maybe show this only with verbose option
		// fmt.Printf("Column X (%d): %v\n", xColumn, xTrimmed)
		// fmt.Printf("Column Y (%v): %v\n", *yColumns, sYTrimmed)
		fmt.Printf("Count: %d, Trim Start: %d, Trim End: %d\n", len(xTrimmed), trimStart, trimEnd)

		regression.PlotTimeData(xTrimmed, sYTrimmed, regression.PlotSettings{
			Title:  pTitle,
			XLabel: pXLabel,
			YLabel: pYLabel,
		})
		err = printCSVColumnStats(remaining, (*yColumns)[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	} else if opt.Called("x") && opt.Called("y") {
		cf := csvutil.New(remaining...)
		cf.NoHeader = noHeader
		cf.FilterZero = filterZero
		query := []int{xColumn}
		query = append(query, (*yColumns)...)
		sliceDatasets, err := cf.GetFloat64Columns(query...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		xTrimmed, err := trimSlice(sliceDatasets[0], trimStart, trimEnd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		var sYTrimmed [][]float64
		for _, ySliceDataset := range sliceDatasets[1:] {
			yTrimmed, _ := trimSlice(ySliceDataset, trimStart, trimEnd)
			sYTrimmed = append(sYTrimmed, yTrimmed)
		}

		// TODO: maybe show this only with verbose option
		fmt.Printf("Column X (%d): %v\n", xColumn, xTrimmed)
		fmt.Printf("Column Y (%v): %v\n", *yColumns, sYTrimmed)
		fmt.Printf("Count: %d, Trim Start: %d, Trim End: %d\n", len(xTrimmed), trimStart, trimEnd)

		regression.PlotRegression(xTrimmed, sYTrimmed, func(x float64) float64 { return x }, 0, regression.PlotSettings{
			Title:     pTitle,
			XLabel:    pXLabel,
			YLabel:    pYLabel,
			DataLabel: pTitle,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		if !opt.Called("regression") {
			os.Exit(0)
		}

		// Original data
		solution, err := regression.SolveTransformation(xTrimmed, sYTrimmed[0], &regression.None{})
		if err == nil {
			err = solution.Plot(&regression.None{})
			printError(err)
		} else {
			printError(err)
		}

		ltList := []interface{}{
			// Exp: y = aB^x | log y = log a + log B * x
			&regression.Exponential{},
			// Power: y = ax^b -> log y = log a + b * log x
			&regression.Power{},
			// y = ln(ax^b) | y = ln a + b * ln x
			&regression.LnPower{},
			// y = 1 / (a + bx) | 1/y = a + bx
			&regression.OneOverX{},
			// y = a + b / (1 + x) | y = a + b * 1/(1+x)
			&regression.BOverX{},
			// y = 1 / (a + bx)^2 | 1/sqrt(y) = a + bx
			&regression.OneOverX2{},
			// y = a + b * sqrt(x) | a + b * sqrt(x)
			&regression.Sqrt{},
		}

		for _, lt := range ltList {
			solution, err = regression.SolveTransformation(
				xTrimmed, sYTrimmed[0], lt.(regression.LinearTransformation))
			if err == nil {
				if review {
					err = solution.PlotLinearTransformation(lt.(regression.Plotter))
					printError(err)
				}
				err = solution.Plot(lt.(regression.Plotter))
				printError(err)
			} else {
				printError(err)
			}
		}

		s, err := regression.SolvePolynomial(xTrimmed, sYTrimmed[0], degree)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		log.Printf("S:\n%3.3g\n", mat64.Formatted(s.A, mat64.Prefix(""), mat64.Squeeze()))

		// si, err := regression.SolvePolynomialReverseMatrix(xSliceDataset, ySliceDataset, degree)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		// 	os.Exit(1)
		// }
		// log.Printf("S (matrix):\n%3.3g\n", mat64.Formatted(si.A, mat64.Prefix(""), mat64.Squeeze()))

		s.Plot()
	} else {
		// Get column stats
		err := printCSVColumnStats(remaining, column)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

}
