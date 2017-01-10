// This file is part of csv-analysis.
//
// Copyright (C) 2017  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package regression

import (
	"io/ioutil"
	"log"
	"math"
	"testing"
)

func testSliceMean(t *testing.T) {
	x := []float64{0, 1, 2, 3, 4}
	xmean := sliceMean(x)
	if xmean != 5 {
		t.Errorf("Slice mean value differs %10g != %f\n", xmean, 5.0)
	}
}

func TestPolynomialRegression(t *testing.T) {
	x := []float64{0, 1, 2, 3, 4}
	y := []float64{0, 1, 2, 3, 4}
	m, err := PolynomialRegression(1, x, y)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	if m.At(0, 0) != 0.0 {
		t.Errorf("Value differs %10g != %f\n", m.At(0, 0), 0.0)
	}
	if m.At(1, 0) != 1.0 {
		t.Errorf("Value differs %10g != %f\n", m.At(1, 0), 1.0)
	}
}

func TestR2Calc(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	x := []float64{0, 1, 2, 3, 4}
	y := []float64{0, 1, 2, 3, 4}
	a := 0.0
	b := 1.0
	s := Solution{At: a, Bt: b}
	r2 := r2Calc(x, y, s.LinearFunction())
	if diff := math.Abs(r2 - 1); diff >= 0.1 {
		t.Errorf("R2 value differs %10g != %f\n", r2, 1.0)
	}
	x = []float64{0, 1, 2, 3, 4}
	y = []float64{0, 0.9, 2, 3.1, 4}
	a = -0.04
	b = 1.02
	s = Solution{At: a, Bt: b}
	r2 = r2Calc(x, y, s.LinearFunction())
	if diff := math.Abs(r2 - 0.9985); diff >= 0.0001 {
		t.Errorf("R2 value differs %10g != %f\n", r2, 0.9985)
	}
}
