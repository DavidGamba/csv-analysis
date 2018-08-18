// This file is part of csv-analysis.
//
// Copyright (C) 2017  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package regression provides linear transformation functions.
*/
package regression

import (
	"fmt"
	"image/color"
	"regexp"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// filenameClean - Cleans the filename
func filenameClean(f string) string {
	f = strings.Replace(f, " ", "_", -1)
	f = strings.Replace(f, "(", "_", -1)
	f = strings.Replace(f, ")", "", -1)
	f = strings.Replace(f, "/", "_over_", -1)
	f = strings.ToLower(f)
	r := regexp.MustCompile(`[^\d\w]`)
	f = r.ReplaceAllString(f, "")
	return f
}

// PlotSettings -
type PlotSettings struct {
	Title     string
	XLabel    string
	YLabel    string
	DataLabel string
	Bold      bool
}

// NewPlot -
func NewPlot(ps PlotSettings) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return p, err
	}
	p.Title.Text = ps.Title
	p.X.Label.Text = ps.XLabel
	p.Y.Label.Text = ps.YLabel
	p.Legend.YOffs = -150
	p.Legend.Top = true
	return p, nil
}

// YDataLabel -
type YDataLabel struct {
	Data  plotter.XYs
	Label string
}

// RGBA color "image/color"
// First color in this list will be used last.
var colorList = []color.RGBA{
	color.RGBA{R: 255, G: 127, B: 127, A: 255}, // Light pink
	color.RGBA{G: 255, A: 255},                 // Green
	color.RGBA{B: 255, A: 255},                 // Blue
	color.RGBA{R: 255, G: 255, A: 255},         // Yellow
	color.RGBA{R: 255, B: 255, A: 255},         // Pink
	color.RGBA{G: 255, B: 255, A: 255},         // Cyan
}

// Gets the next color available.
// When i=0, it will return RED, otherwise, it will return one color from the
// color list.
func getColor(i int) color.RGBA {
	if i == 0 {
		return color.RGBA{R: 255, A: 255} // Red
	}
	n := len(colorList)
	j := i % n
	return colorList[j]
}

// PlotRegression -
func PlotRegression(x []float64, ys [][]float64, f func(float64) float64, r2 float64, ps PlotSettings) error {
	p, err := NewPlot(ps)
	if err != nil {
		return err
	}
	// Plot in reverse order because the first entry is the most important
	last := len(ys) - 1
	for index := range ys {
		i := last - index
		y := ys[i]
		pts := make(plotter.XYs, len(x))
		for j := range x {
			pts[j].X = x[j]
			pts[j].Y = y[j]
		}
		lpLine, lpPoints, err := plotter.NewLinePoints(pts)
		if err != nil {
			return err
		}
		lpLine.Color = getColor(i)
		lpPoints.Color = getColor(i)
		p.Add(lpLine, lpPoints)
		p.Legend.Add(fmt.Sprintf("%s %d", ps.DataLabel, i), lpLine, lpPoints)
	}

	if r2 != 0 {
		pf := plotter.NewFunction(f)
		p.Add(pf)
		p.Legend.Add("Regression", pf)
		p.Legend.Add(fmt.Sprintf("R² %.4f", r2))
		// p.Legend.Add(fmt.Sprintf("y = %7f + %7fx", s.At(0, 0), s.At(1, 0)))
	}

	name := "plot-" + filenameClean(ps.Title) + ".png"
	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, name); err != nil {
		return err
	}
	fmt.Printf("Plot: %s\n", name)

	return nil
}

// PlotTimeData -
func PlotTimeData(x []float64, ys [][]float64, ps PlotSettings) error {
	p, err := NewPlot(ps)
	if err != nil {
		return err
	}
	p.X.Tick.Marker = plot.TimeTicks{}
	// Plot in reverse order because the first entry is the most important
	last := len(ys) - 1
	for index := range ys {
		i := last - index
		y := ys[i]
		pts := make(plotter.XYs, len(x))
		for j := range x {
			pts[j].X = x[j]
			pts[j].Y = y[j]
		}
		lpLine, _, err := plotter.NewLinePoints(pts)
		if err != nil {
			return err
		}
		lpLine.Color = getColor(i)
		if i == 0 && ps.Bold {
			lpLine.Width = 2
		}
		p.Add(lpLine)
		p.Legend.Add(fmt.Sprintf("%s %d", ps.DataLabel, i), lpLine)
	}
	// Save the plot to a PNG file.
	name := "plot-" + filenameClean(ps.Title) + ".png"
	// 6, 3.5
	if err := p.Save(12*vg.Inch, 7*vg.Inch, name); err != nil {
		return err
	}
	fmt.Printf("Plot: %s\n", name)

	return nil
}

// PlotLinearTransformation -
func (s Solution) PlotLinearTransformation(p Plotter) error {
	fmt.Printf("Linear   %-20s R²=%.4f a=%10f b=%10f\n", p.Name(), s.R2t, s.At, s.Bt)
	fmt.Printf("         %s -> %s\n", p.TextEquation(), p.TextTransformedEquation())
	return PlotRegression(s.Xt, [][]float64{s.Yt}, s.LinearFunction(), s.R2t, PlotSettings{
		Title:     "Linear " + p.Name(),
		XLabel:    p.XLabel(),
		YLabel:    p.YLabel(),
		DataLabel: "Data",
	})
}

// Plot -
func (s Solution) Plot(p Plotter) error {
	fmt.Printf("Equation %-20s R²t=%.4f R²=%.4f a=%10f b=%10f\n", p.TextEquation(), s.R2t, s.R2, s.A, s.B)
	return PlotRegression(s.X, [][]float64{s.Y}, s.RegressionFunction(), s.R2, PlotSettings{
		Title:     p.Name(),
		XLabel:    "X",
		YLabel:    "Y",
		DataLabel: "Data",
	})
}

// Plot -
func (s PolynomialSolution) Plot() error {
	fmt.Printf("Polynomial degree %d R²=%.4f\n", s.Degree, s.R2)
	return PlotRegression(s.X, [][]float64{s.Y}, s.PolynomialFunction(), s.R2, PlotSettings{
		Title:     "Polynomial Regression",
		XLabel:    "X",
		YLabel:    "Y",
		DataLabel: "Data",
	})
}
