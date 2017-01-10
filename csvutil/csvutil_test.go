// This file is part of csv-analysis.
//
// Copyright (C) 2017  David Gamba Rios
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package csvutil

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetCSVColumns(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	in := `a,b,c
1,2,3
1,2,3
1,2,3
`
	_, err := getCSVColumns(strings.NewReader(in), 0)
	if err == nil {
		t.Fatalf("Expected error not thrown\n")
	}
	if err.Error() != "Column index error: 0 <= 0!" {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	cdata, err := getCSVColumns(strings.NewReader(in), 1)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	expected := [][]string{[]string{"a", "1", "1", "1"}}
	if !reflect.DeepEqual(cdata, expected) {
		t.Errorf("Wrong data: %v != %v\n", cdata, expected)
	}
	cdata, err = getCSVColumns(strings.NewReader(in), 2)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	expected = [][]string{[]string{"b", "2", "2", "2"}}
	if !reflect.DeepEqual(cdata, expected) {
		t.Errorf("Wrong data: %v != %v\n", cdata, expected)
	}
	cdata, err = getCSVColumns(strings.NewReader(in), 3)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	expected = [][]string{[]string{"c", "3", "3", "3"}}
	if !reflect.DeepEqual(cdata, expected) {
		t.Errorf("Wrong data: %v != %v\n", cdata, expected)
	}
	cdata, err = getCSVColumns(strings.NewReader(in), 4)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	expected = make([][]string, 1)
	if !reflect.DeepEqual(cdata, expected) {
		t.Errorf("Wrong data: %v != %v\n", cdata, expected)
	}
}

func TestGetCSVRows(t *testing.T) {
	log.SetOutput(os.Stderr)
	in := `a,b,c
1,1,1
2,2,2
`
	_, err := getCSVRows(strings.NewReader(in), 0)
	if err == nil {
		t.Fatalf("Expected error not thrown\n")
	}
	if err.Error() != "Row index error: 0 <= 0!" {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	rdata, err := getCSVRows(strings.NewReader(in), 1)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	expected := [][]string{[]string{"a", "b", "c"}}
	if !reflect.DeepEqual(rdata, expected) {
		t.Errorf("Wrong data: %v != %v\n", rdata, expected)
	}
	rdata, err = getCSVRows(strings.NewReader(in), 2)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	expected = [][]string{[]string{"1", "1", "1"}}
	if !reflect.DeepEqual(rdata, expected) {
		t.Errorf("Wrong data: %v != %v\n", rdata, expected)
	}
	rdata, err = getCSVRows(strings.NewReader(in), 3)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	expected = [][]string{[]string{"2", "2", "2"}}
	if !reflect.DeepEqual(rdata, expected) {
		t.Errorf("Wrong data: %v != %v\n", rdata, expected)
	}
	rdata, err = getCSVRows(strings.NewReader(in), 4)
	if err != nil {
		t.Fatalf("Unexpected error: %s\n", err)
	}
	expected = make([][]string, 1)
	if !reflect.DeepEqual(rdata, expected) {
		t.Errorf("Wrong data: %v != %v\n", rdata, expected)
	}
}
