// This file is part of csv-analysis.
//
// Copyright (C) 2017  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package stat provides ways to analyse data from one or more files and generate basic statistic analysis.
*/
package stat

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"os"
)

// PrintSliceStats -
func PrintSliceStats(data []float64) {
	fmt.Printf("Count: %d\n", len(data))
	var d stats.Float64Data = data
	max, err := d.Max()
	printError(err)
	fmt.Printf("Max: %f\n", max)
	min, err := d.Min()
	printError(err)
	fmt.Printf("Min: %f\n", min)
	mean, err := d.Mean()
	printError(err)
	fmt.Printf("Mean: %f\n", mean)
	sd, err := d.StandardDeviation()
	printError(err)
	// TODO: Get better errors
	fmt.Printf("Standard Deviation σ: %f, %f%%\n", sd, sd*100/mean)
	variance, err := d.Variance()
	printError(err)
	fmt.Printf("Variance σ²: %f\n", variance)
	median, err := d.Median()
	printError(err)
	fmt.Printf("Median: %f\n", median)
	medianDeviation, err := d.MedianAbsoluteDeviation()
	printError(err)
	fmt.Printf("Median Absolute Deviation MAD: %f, %f%%\n", medianDeviation, medianDeviation*100/median)
	sum, err := d.Sum()
	printError(err)
	fmt.Printf("Sum: %f\n", sum)
}

// printError - prints the given error to STDERR.
func printError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
}
