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

// trimData -
func trimData(x, y []float64, trimStart, trimEnd int) ([]float64, []float64) {
	n := len(x) - trimStart - trimEnd
	xSliceDataset := make([]float64, n)
	ySliceDataset := make([]float64, n)
	for i := range x {
		if i < trimStart {
			continue
		}
		if i >= len(x)-trimEnd {
			break
		}
		xSliceDataset[i-trimStart] = x[i]
		ySliceDataset[i-trimStart] = y[i]
	}
	return xSliceDataset, ySliceDataset
}

func synopsis() {
	synopsis := `script --column|-c <n> <csv-file>...
       [--no-header|--nh] [--filter-zero|--fz] 

# Regression analysis
script -x <n> -y <n> <csv-file>...
       [--no-header|--nh] [--filter-zero|--fz]
			 [--trim-start|--ts <n>] [--trim-end|--te <n>]
			 [--degree] [--review]

# Time plot
script -x <n> -y <n> <csv-file>... -xtime <timeformat>
       [--no-header|--nh] [--filter-zero|--fz]
			 [--trim-start|--ts <n>] [--trim-end|--te <n>]

# Inspect data and exit
script [--show-header|-s] [--show-data|--sd] <csv-file>...

script [--help]

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
`
	fmt.Fprintln(os.Stderr, synopsis)
}

func main() {
	var column, xColumn, yColumn int // field to analize
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
	opt.IntVar(&yColumn, "y", 1)
	// CSV data trimming
	opt.IntVar(&trimStart, "trim-start", 0, "ts")
	opt.IntVar(&trimEnd, "trim-end", 0, "te")
	// Linear Regression degree
	opt.IntVar(&degree, "degree", 1, "degree")
	// Action
	opt.Bool("plot-data", false, "p")
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
		var xSliceDataset, ySliceDataset []float64
		for _, e := range sliceDatasetsString[0] {
			trimmed := strings.TrimSpace(e)
			t, err := time.Parse(xTimeFormat, trimmed)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: time format '%s': %s\n", xTimeFormat, err)
				continue
			}
			xSliceDataset = append(xSliceDataset, float64(t.Unix()))
		}
		sliceDatasets, err := cf.GetFloat64Columns(yColumn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		ySliceDataset = sliceDatasets[0]
		// TODO: maybe show this only with verbose option
		// fmt.Printf("Column X (%d): %v\n", xColumn, xSliceDataset)
		fmt.Printf("Column Y (%d): %v\n", yColumn, ySliceDataset)
		fmt.Printf("Count: %d, Trim Start: %d, Trim End: %d\n", len(xSliceDataset), trimStart, trimEnd)

		if len(xSliceDataset) != len(ySliceDataset) {
			fmt.Fprintf(os.Stderr, "ERROR: Column lenghts do not match\n")
			os.Exit(1)
		}
		xTrimmed, yTrimmed := trimData(xSliceDataset, ySliceDataset, trimStart, trimEnd)
		regression.PlotTimeData(xTrimmed, yTrimmed, regression.PlotSettings{
			Title:  pTitle,
			XLabel: pXLabel,
			YLabel: pYLabel,
		})
		err = printCSVColumnStats(remaining, yColumn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	} else if opt.Called("x") && opt.Called("y") {
		cf := csvutil.New(remaining...)
		cf.NoHeader = noHeader
		cf.FilterZero = filterZero
		sliceDatasets, err := cf.GetFloat64Columns(xColumn, yColumn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		xSliceDataset := sliceDatasets[0]
		ySliceDataset := sliceDatasets[1]

		// TODO: maybe show this only with verbose option
		fmt.Printf("Column X (%d): %v\n", xColumn, xSliceDataset)
		fmt.Printf("Column Y (%d): %v\n", yColumn, ySliceDataset)
		fmt.Printf("Count: %d, Trim Start: %d, Trim End: %d\n", len(xSliceDataset), trimStart, trimEnd)

		xTrimmed, yTrimmed := trimData(xSliceDataset, ySliceDataset, trimStart, trimEnd)
		if opt.Called("plot-data") {
			regression.PlotRegression(xTrimmed, yTrimmed, func(x float64) float64 { return x }, 0, regression.PlotSettings{Title: "Data", XLabel: "X", YLabel: "Y", DataLabel: "Data"})
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}

		// Original data
		solution, err := regression.SolveTransformation(xTrimmed, yTrimmed, &regression.None{})
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
				xTrimmed, yTrimmed, lt.(regression.LinearTransformation))
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

		s, err := regression.SolvePolynomial(xTrimmed, yTrimmed, degree)
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
