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
	// "log"
	"math"
)

// LinearTransformation -
type LinearTransformation interface {
	Name() string
	FX(a, b, x float64) float64    // fx(x)
	FTransformX(x float64) float64 // function to transform X
	FTransformY(y float64) float64 // function to transform Y
	FRestoreA(at float64) float64  // function to restore A
	FRestoreB(bt float64) float64  // function to restore B
}

// Interpolation - Allows to get an X or Y point based on y or x.
type Interpolation interface {
	FY(a, b, y float64) float64 // fy(y) solve the equation for y, used for interpolation.
	FX(a, b, x float64) float64 // fx(x)
}

// Plotter -
type Plotter interface {
	Name() string
	TextEquation() string
	TextTransformedEquation() string
	TLabel() string             // Top label
	XLabel() string             // X label
	YLabel() string             // Y label
	FX(a, b, x float64) float64 // fx(x)
}

// None - No transformation
type None struct{}

// Name - Name
func (*None) Name() string { return "No Transformation" }

// TextEquation - ASCII equation
func (*None) TextEquation() string { return "y = a + bx" }

// TextTransformedEquation - Text equation
func (*None) TextTransformedEquation() string { return "y = a + bx" }

// TLabel - Top label
func (*None) TLabel() string { return "y vs x" }

// XLabel - X label
func (*None) XLabel() string { return "x" }

// YLabel - Y label
func (*None) YLabel() string { return "y" }

// FX - Function to calculate normal y values
//    y = a + bx
func (*None) FX(a, b, x float64) (y float64) {
	return a + b*x
}

// FY - Solve the equation for y, used for interpolation
//    x = (y - a)/b
func (*None) FY(a, b, y float64) (x float64) {
	return (y - a) / b
}

// FTransformX - function to transform X
func (*None) FTransformX(x float64) (xt float64) {
	return x
}

// FTransformY - function to transform Y
func (*None) FTransformY(y float64) (yt float64) {
	return y
}

// FRestoreA - function to restore A
func (*None) FRestoreA(at float64) (a float64) {
	return at
}

// FRestoreB - function to restore B
func (*None) FRestoreB(bt float64) (b float64) {
	return bt
}

// Power - Provides the transformation functions that satisfy:
//    y = ax^b
// to:
//    log y = log a + b * log x
type Power struct{}

// Name - Name
func (*Power) Name() string { return "Power" }

// TextEquation - Text equation
func (*Power) TextEquation() string { return "y = ax^b" }

// TextTransformedEquation - Text equation
func (*Power) TextTransformedEquation() string { return "log y = log a + b * log x" }

// TLabel - Top label
func (*Power) TLabel() string { return "log(y) vs log(x)" }

// XLabel - X label
func (*Power) XLabel() string { return "log(x)" }

// YLabel - Y label
func (*Power) YLabel() string { return "log(y)" }

// FX - Function to calculate normal y values
//    y = ax^b
func (*Power) FX(a, b, x float64) (y float64) {
	// log.Printf("a: %f, b: %f, x: %f\n", a, b, x)
	return a * math.Pow(x, b)
}

// FY - Solve the equation for y, used for interpolation
//    x = sqrt(y/a,b)
func (*Power) FY(a, b, y float64) (x float64) {
	// log.Printf("a: %f, b: %f, x: %f\n", a, b, x)
	return a * math.Pow(x, b)
}

// FTransformX - function to transform X
func (*Power) FTransformX(x float64) (xt float64) {
	return math.Log10(x)
}

// FTransformY - function to transform Y
func (*Power) FTransformY(y float64) (yt float64) {
	return math.Log10(y)
}

// FRestoreA - function to restore A
func (*Power) FRestoreA(at float64) (a float64) {
	return math.Pow(10, at)
}

// FRestoreB - function to restore B
func (*Power) FRestoreB(bt float64) (b float64) {
	return bt
}

// Exponential - Provides the transformation functions that satisfy:
//    y = ae^(bx) = aB^x
// to:
//    ln y = ln a + bx = ln a + ln B * x
type Exponential struct{}

// Name - Name
func (*Exponential) Name() string { return "Exponential" }

// TextEquation - ASCII equation
func (*Exponential) TextEquation() string { return "y = aB^x" }

// TextTransformedEquation - Text equation
func (*Exponential) TextTransformedEquation() string { return "ln y = ln a + ln B * x" }

// TLabel - Top label
func (*Exponential) TLabel() string { return "ln(y) vs x" }

// XLabel - X label
func (*Exponential) XLabel() string { return "x" }

// YLabel - Y label
func (*Exponential) YLabel() string { return "ln(y)" }

// FX - Function to calculate normal y values
//    y = aB^x
func (*Exponential) FX(a, b, x float64) (y float64) {
	return a * math.Pow(b, x)
}

// FY - Solve the equation for y, used for interpolation
//    x = ln(y/a)/ln(b)
func (*Exponential) FY(a, b, y float64) (x float64) {
	return math.Log(y/a) / math.Log(b)
}

// FTransformX - function to transform X
func (*Exponential) FTransformX(x float64) (xt float64) {
	return x
}

// FTransformY - function to transform Y
func (*Exponential) FTransformY(y float64) (yt float64) {
	return math.Log(y)
}

// FRestoreA - function to restore A
func (*Exponential) FRestoreA(at float64) (a float64) {
	return math.Pow(math.E, at)
}

// FRestoreB - function to restore B
func (*Exponential) FRestoreB(bt float64) (b float64) {
	return math.Pow(math.E, bt)
}

// LnPower - Provides the transformation functions that satisfy:
//    y = ln(ax^b)
// to:
//    y = ln a + b * ln x
type LnPower struct{}

// Name - Name
func (*LnPower) Name() string { return "LnPower" }

// TextEquation - ASCII equation
func (*LnPower) TextEquation() string { return "y = ln(ax^b)" }

// TextTransformedEquation - Text equation
func (*LnPower) TextTransformedEquation() string { return "y = ln a + b ln x" }

// TLabel - Top label
func (*LnPower) TLabel() string { return "y vs ln(x)" }

// XLabel - X label
func (*LnPower) XLabel() string { return "ln(x)" }

// YLabel - Y label
func (*LnPower) YLabel() string { return "y" }

// FX - Function to calculate normal y values
//    y = ln a + b ln x
func (*LnPower) FX(a, b, x float64) (y float64) {
	return math.Log(a) + b*math.Log(x)
}

// FY - Solve the equation for y, used for interpolation
//    x = ((y - b)/ln a)^e
func (*LnPower) FY(a, b, y float64) (x float64) {
	return math.Pow((y-b)/math.Log(a), math.E)
}

// FTransformX - function to transform X
func (*LnPower) FTransformX(x float64) (xt float64) {
	return math.Log(x)
}

// FTransformY - function to transform Y
func (*LnPower) FTransformY(y float64) (yt float64) {
	return y
}

// FRestoreA - function to restore A
func (*LnPower) FRestoreA(at float64) (a float64) {
	return math.Pow(math.E, at)
}

// FRestoreB - function to restore B
func (*LnPower) FRestoreB(bt float64) (b float64) {
	return bt
}

// OneOverX - Provides the transformation functions that satisfy:
//    y = 1 / (a + bx)
// to:
//    1/y = a + bx
type OneOverX struct{}

// Name - Name
func (*OneOverX) Name() string { return "OneOverX" }

// TextEquation - ASCII equation
func (*OneOverX) TextEquation() string { return "y = 1 / (a + bx)" }

// TextTransformedEquation - Text equation
func (*OneOverX) TextTransformedEquation() string { return "1/y = a + bx" }

// TLabel - Top label
func (*OneOverX) TLabel() string { return "1/y vs x" }

// XLabel - X label
func (*OneOverX) XLabel() string { return "x" }

// YLabel - Y label
func (*OneOverX) YLabel() string { return "1/y" }

// FX - Function to calculate normal y values
//    y = 1 / (a + bx)
func (*OneOverX) FX(a, b, x float64) (y float64) {
	return 1 / (a + b*x)
}

// FY - Solve the equation for y, used for interpolation
//    x = (1/y - a) / b
func (*OneOverX) FY(a, b, y float64) (x float64) {
	return (1/y - a) / b
}

// FTransformX - function to transform X
func (*OneOverX) FTransformX(x float64) (xt float64) {
	return x
}

// FTransformY - function to transform Y
func (*OneOverX) FTransformY(y float64) (yt float64) {
	return 1 / y
}

// FRestoreA - function to restore A
func (*OneOverX) FRestoreA(at float64) (a float64) {
	return at
}

// FRestoreB - function to restore B
func (*OneOverX) FRestoreB(bt float64) (b float64) {
	return bt
}

// BOverX - Provides the transformation functions that satisfy:
//    y = a + b / (1 + x)
// to:
//    y = a + b * 1 / (1 + x)
type BOverX struct{}

// Name - Name
func (*BOverX) Name() string { return "BOverX" }

// TextEquation - ASCII equation
func (*BOverX) TextEquation() string { return "y = a + b / (1 + x)" }

// TextTransformedEquation - Text equation
func (*BOverX) TextTransformedEquation() string { return "y = a + b *  1 / (1 + x)" }

// TLabel - Top label
func (*BOverX) TLabel() string { return "y vs 1/(1+x)" }

// XLabel - X label
func (*BOverX) XLabel() string { return "1/(1+x)" }

// YLabel - Y label
func (*BOverX) YLabel() string { return "y" }

// FX - Function to calculate normal y values
//    y = 1 / (a + bx)
func (*BOverX) FX(a, b, x float64) (y float64) {
	return a + b/(1+x)
}

// FY - Solve the equation for y, used for interpolation
//    x = (1/((y -a ) / b)) - 1
func (*BOverX) FY(a, b, y float64) (x float64) {
	return (1 / ((y - a) / b)) - 1
}

// FTransformX - function to transform X
func (*BOverX) FTransformX(x float64) (xt float64) {
	return 1 / (1 + x)
}

// FTransformY - function to transform Y
func (*BOverX) FTransformY(y float64) (yt float64) {
	return y
}

// FRestoreA - function to restore A
func (*BOverX) FRestoreA(at float64) (a float64) {
	return at
}

// FRestoreB - function to restore B
func (*BOverX) FRestoreB(bt float64) (b float64) {
	return bt
}

// OneOverX2 - Provides the transformation functions that satisfy:
//    y = 1 / (a + bx)^2
// to:
//    1/sqrt(y) = a + bx
type OneOverX2 struct{}

// Name - Name
func (*OneOverX2) Name() string { return "OneOverX2" }

// TextEquation - ASCII equation
func (*OneOverX2) TextEquation() string { return "y = 1 / (a + bx)^2" }

// TextTransformedEquation - Text equation
func (*OneOverX2) TextTransformedEquation() string { return "1/sqrt(y) = a + bx" }

// TLabel - Top label
func (*OneOverX2) TLabel() string { return "1/sqrt(y) vs x" }

// XLabel - X label
func (*OneOverX2) XLabel() string { return "x" }

// YLabel - Y label
func (*OneOverX2) YLabel() string { return "1/sqrt(y)" }

// FX - Function to calculate normal y values
//    y = 1 / (a + bx)^2
func (*OneOverX2) FX(a, b, x float64) (y float64) {
	return 1 / math.Pow(a+b*x, float64(2))
}

// FY - Solve the equation for y, used for interpolation
//    x = (1/sqrt(y) - a)/b
func (*OneOverX2) FY(a, b, y float64) (x float64) {
	return (1/math.Sqrt(y) - a) / b
}

// FTransformX - function to transform X
func (*OneOverX2) FTransformX(x float64) (xt float64) {
	return x
}

// FTransformY - function to transform Y
func (*OneOverX2) FTransformY(y float64) (yt float64) {
	return 1 / math.Sqrt(y)
}

// FRestoreA - function to restore A
func (*OneOverX2) FRestoreA(at float64) (a float64) {
	return at
}

// FRestoreB - function to restore B
func (*OneOverX2) FRestoreB(bt float64) (b float64) {
	return bt
}

// Sqrt - Provides the transformation functions that satisfy:
//    y = a + b * sqrt(x)
// to:
//    y = a + b * sqrt(x)
type Sqrt struct{}

// Name - Name
func (*Sqrt) Name() string { return "Sqrt" }

// TextEquation - ASCII equation
func (*Sqrt) TextEquation() string { return "y = a + b * sqrt(x)" }

// TextTransformedEquation - Text equation
func (*Sqrt) TextTransformedEquation() string { return "y = a + b * sqrt(x)" }

// TLabel - Top label
func (*Sqrt) TLabel() string { return "y vs sqrt(x)" }

// XLabel - X label
func (*Sqrt) XLabel() string { return "sqrt(x)" }

// YLabel - Y label
func (*Sqrt) YLabel() string { return "y" }

// FX - Function to calculate normal y values
//    y = a + b * sqrt(x)
func (*Sqrt) FX(a, b, x float64) (y float64) {
	return a + b*math.Sqrt(x)
}

// FY - Solve the equation for y, used for interpolation
//    x = ((y - a)/b)^2
func (*Sqrt) FY(a, b, y float64) (x float64) {
	return math.Pow((y-a)/b, float64(2))
}

// FTransformX - function to transform X
func (*Sqrt) FTransformX(x float64) (xt float64) {
	return math.Sqrt(x)
}

// FTransformY - function to transform Y
func (*Sqrt) FTransformY(y float64) (yt float64) {
	return y
}

// FRestoreA - function to restore A
func (*Sqrt) FRestoreA(at float64) (a float64) {
	return at
}

// FRestoreB - function to restore B
func (*Sqrt) FRestoreB(bt float64) (b float64) {
	return bt
}
