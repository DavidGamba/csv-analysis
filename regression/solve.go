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
	"log"
	"math"

	"gonum.org/v1/gonum/mat"
)

// Solution - Linear Transformation Solution.
type Solution struct {
	X, Y   []float64 // original data slices.
	Xt, Yt []float64 // transformed data slices
	At, Bt float64   // solution to the transformed linear equation
	A, B   float64   // solution to the linear equation
	LT     LinearTransformation
	R2t    float64
	R2     float64
}

// SolveTransformation - Given a pointer to X and Y []float64 data and a linear
// transformation function, it will return the solution.
func SolveTransformation(xo, yo []float64, lt LinearTransformation) (Solution, error) {
	result := Solution{LT: lt}
	n := len(xo)
	result.X = make([]float64, n)
	result.Y = make([]float64, n)
	result.Xt = make([]float64, n)
	result.Yt = make([]float64, n)
	for i := range xo {
		result.X[i] = xo[i]
		result.Y[i] = yo[i]
		result.Xt[i] = lt.FTransformX(xo[i])
		result.Yt[i] = lt.FTransformY(yo[i])
	}

	s, err := PolynomialRegression(1, result.Xt, result.Yt)
	if err != nil {
		return result, err
	}
	result.At = s.At(0, 0)
	result.Bt = s.At(1, 0)
	result.A = lt.FRestoreA(result.At)
	result.B = lt.FRestoreB(result.Bt)

	result.R2t = r2Calc(result.Xt, result.Yt, result.LinearFunction())
	result.R2 = r2Calc(result.X, result.Y, result.RegressionFunction())
	return result, nil
}

// LinearFunction - Returns a linear function based on the transformed At and Bt values.
func (s Solution) LinearFunction() func(x float64) float64 {
	return func(x float64) float64 {
		y := s.At + s.Bt*x
		return y
	}
}

// RegressionFunction -
func (s Solution) RegressionFunction() func(x float64) float64 {
	return func(x float64) float64 {
		return s.LT.FX(s.A, s.B, x)
	}
}

// PolynomialSolution - Polynomial solution.
type PolynomialSolution struct {
	X, Y   []float64 // Original data slices.
	Degree int
	A      *mat.Dense // Solution
	R2     float64
}

// SolvePolynomial - Given a pointer to X and Y []float64 data and a linear
// transformation function, it will return the solution.
func SolvePolynomial(xo, yo []float64, degree int) (PolynomialSolution, error) {
	result := PolynomialSolution{Degree: degree}
	n := len(xo)
	result.X = make([]float64, n)
	result.Y = make([]float64, n)
	for i := range xo {
		result.X[i] = xo[i]
		result.Y[i] = yo[i]
	}

	log.Printf("Polynomial regression of degree %d\n", degree)
	s, err := PolynomialRegression(result.Degree, result.X, result.Y)
	if err != nil {
		return result, err
	}
	log.Printf("S:\n%3.3g\n", mat.Formatted(s, mat.Prefix(""), mat.Squeeze()))

	result.A = s
	result.R2 = r2Calc(result.X, result.Y, result.PolynomialFunction())
	return result, nil
}

// SolvePolynomialReverseMatrix - Given a pointer to X and Y []float64 data and a linear
// transformation function, it will return the solution using the Reverse Matrix.
func SolvePolynomialReverseMatrix(xo, yo []float64, degree int) (PolynomialSolution, error) {
	result := PolynomialSolution{Degree: degree}
	n := len(xo)
	result.X = make([]float64, n)
	result.Y = make([]float64, n)
	for i := range xo {
		result.X[i] = xo[i]
		result.Y[i] = yo[i]
	}

	log.Printf("Polynomial regression of degree %d\n", degree)
	s, err := polynomialRegressionInverseMatrix(result.Degree, result.X, result.Y)
	if err != nil {
		return result, err
	}
	log.Printf("S:\n%3.3g\n", mat.Formatted(s, mat.Prefix(""), mat.Squeeze()))

	result.A = s.(*mat.Dense)
	result.R2 = r2Calc(result.X, result.Y, result.PolynomialFunction())
	return result, nil
}

// PolynomialFunction - Returns a linear function based on the provided A matrix.
func (s PolynomialSolution) PolynomialFunction() func(x float64) float64 {
	r, _ := s.A.Caps()
	// if c != 1 {
	// 	return nil, fmt.Errorf("Unexpected number of rows\n")
	// }
	return func(x float64) float64 {
		var y float64
		for i := 0; i < r; i++ {
			y += s.A.At(i, 0) * math.Pow(x, float64(i))
		}
		return y
	}
}

// PolynomialRegression - Least Squares Polynomial Regression
// Returns A and B matrix based on the given polynomial degree.
// m: polynomial degree
// n: number of points
// When m = 2:
//                                  A = B
//        (n)a₀ +  (∑xᵢ)a₁ + (∑xᵢ²)a₂ = ∑yᵢ
//      (∑xᵢ)a₀ + (∑xᵢ²)a₁ + (∑xᵢ³)a₂ = ∑xᵢyᵢ
//     (∑xᵢ²)a₀ + (∑xᵢ³)a₁ + (∑xᵢ⁴)a₂ = ∑xᵢ²yᵢ
func PolynomialRegression(m int, x, y []float64) (*mat.Dense, error) {
	solved := mat.NewDense(m+1, 1, nil)
	n := len(x)
	if n < m+1 {
		return solved, fmt.Errorf("Not enough points")
	}

	// This is what you would use if you where using slice of slices instead of gonum matrix
	// // Initialize the slice of slices
	// initSliceOfSlices := func(r, c int) [][]float64 {
	// 	a := make([][]float64, r)
	// 	for i := range a {
	// 		a[i] = make([]float64, c)
	// 	}
	// 	return a
	// }

	a := mat.NewDense(m+1, m+1, nil)
	b := mat.NewDense(m+1, 1, nil)
	// as := initSliceOfSlices(m+1, m+1)
	// bs := initSliceOfSlices(m+1, 1)
	for i := 0; i < m+1; i++ {
		for j := 0; j <= i; j++ {
			k := i + j
			// log.Printf("j: %d, k: %d", j, k)
			var sum float64
			sum = 0
			for l := 0; l < n; l++ {
				// s := math.Pow(x[l], float64(k))
				// log.Printf("l: %d, s: %f", l, s)
				sum = sum + math.Pow(x[l], float64(k))
			}
			// log.Printf("a(%d)(%d) sum: %f\n", i, j, sum)
			// log.Printf("a(%d)(%d) sum: %f\n", j, i, sum)
			a.Set(i, j, sum)
			a.Set(j, i, sum)
			// as[i][j] = sum
			// as[j][i] = sum
			// log.Printf("a: %v\n", a)
			// log.Printf("am:\n%v\n", mat.Formatted(am, mat.Prefix("")))
		}
		var sum float64
		sum = 0
		for l := 0; l < n; l++ {
			// s := y[l] * math.Pow(x[l], float64(i))
			// log.Printf("l %d, i %d, s %f\n", l, i, s)
			sum = sum + y[l]*math.Pow(x[l], float64(i))
		}
		b.Set(i, 0, sum)
		// bs[i][0] = sum
		// log.Printf("b:\n%v\n", mat.Formatted(b, mat.Prefix("")))
	}
	// log.Printf("A:\n%3.3g\n", mat.Formatted(a, mat.Prefix(""), mat.Squeeze()))
	// log.Printf("B:\n%3.3g\n", mat.Formatted(b, mat.Prefix(""), mat.Squeeze()))

	solved.Solve(a, b)
	return solved, nil
}

func generateZmatrix(x []float64, m int) *mat.Dense {
	z := mat.NewDense(len(x), m+1, nil)
	for i, xi := range x {
		for j := 0; j <= m; j++ {
			z.Set(i, j, math.Pow(xi, float64(j)))
		}
	}
	return z
}

func polynomialRegressionInverseMatrix(m int, xSliceDataset, ySliceDataset []float64) (mat.Matrix, error) {
	var afinal, zAndzTranspose, zAndzTransposeInverse, yAndzTranspose mat.Dense
	y := mat.NewVecDense(len(ySliceDataset), ySliceDataset)
	z := generateZmatrix(xSliceDataset, m)
	log.Printf("[Z]:\n%f\n", mat.Formatted(z, mat.Prefix(""), mat.Squeeze()))
	zAndzTranspose.Mul(z.T(), z)
	log.Printf("[Z]^T [Z]:\n%f\n", mat.Formatted(&zAndzTranspose, mat.Prefix(""), mat.Squeeze()))
	yAndzTranspose.Mul(z.T(), y)
	log.Printf("[Z]^T {y}:\n%f\n", mat.Formatted(&yAndzTranspose, mat.Prefix(""), mat.Squeeze()))
	zAndzTransposeInverse.Inverse(&zAndzTranspose)
	log.Printf("[[Z]^T [Z]]^1:\n%f\n", mat.Formatted(&zAndzTransposeInverse, mat.Prefix(""), mat.Squeeze()))
	afinal.Mul(&zAndzTransposeInverse, &yAndzTranspose)
	log.Printf("A:\n%f\n", mat.Formatted(&afinal, mat.Prefix(""), mat.Squeeze()))
	return &afinal, nil
}

// sliceMean - Returns 1/n * ∑xᵢ
func sliceMean(x []float64) float64 {
	var mean float64
	n := len(x)
	for _, e := range x {
		mean += e
	}
	return mean / float64(n)
}

// r2Calc - returns the R²
// ∑(yᵢFitted - ySampleMean)² / ∑(yᵢ - ySampleMean)²
func r2Calc(x, y []float64, fx func(float64) float64) float64 {
	mean := sliceMean(y)
	var ssTotal, ssReg float64
	for i, xi := range x {
		ssReg += math.Pow(fx(xi)-mean, float64(2))
		ssTotal += math.Pow(y[i]-mean, float64(2))
	}
	return ssReg / ssTotal
}
