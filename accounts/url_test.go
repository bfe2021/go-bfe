// Copyright 2018 The go-bfe Authors
// This file is part of the go-bfe library.
//
// The go-bfe library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-bfe library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-bfe library. If not, see <http://www.gnu.org/licenses/>.

package accounts

import (
	"testing"
)

func TestURLParsing(t *testing.T) {
	url, err := parseURL("https://bfedu.io")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "bfedu.io" {
		t.Errorf("expected: %v, got: %v", "bfedu.io", url.Path)
	}

	_, err = parseURL("bfedu.io")
	if err == nil {
		t.Error("expected err, got: nil")
	}
}

func TestURLString(t *testing.T) {
	url := URL{Scheme: "https", Path: "bfedu.io"}
	if url.String() != "https://bfedu.io" {
		t.Errorf("expected: %v, got: %v", "https://bfedu.io", url.String())
	}

	url = URL{Scheme: "", Path: "bfedu.io"}
	if url.String() != "bfedu.io" {
		t.Errorf("expected: %v, got: %v", "bfedu.io", url.String())
	}
}

func TestURLMarshalJSON(t *testing.T) {
	url := URL{Scheme: "https", Path: "bfedu.io"}
	json, err := url.MarshalJSON()
	if err != nil {
		t.Errorf("unexpcted error: %v", err)
	}
	if string(json) != "\"https://bfedu.io\"" {
		t.Errorf("expected: %v, got: %v", "\"https://bfedu.io\"", string(json))
	}
}

func TestURLUnmarshalJSON(t *testing.T) {
	url := &URL{}
	err := url.UnmarshalJSON([]byte("\"https://bfedu.io\""))
	if err != nil {
		t.Errorf("unexpcted error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "bfedu.io" {
		t.Errorf("expected: %v, got: %v", "https", url.Path)
	}
}

func TestURLComparison(t *testing.T) {
	tests := []struct {
		urlA   URL
		urlB   URL
		expect int
	}{
		{URL{"https", "bfedu.io"}, URL{"https", "bfedu.io"}, 0},
		{URL{"http", "bfedu.io"}, URL{"https", "bfedu.io"}, -1},
		{URL{"https", "bfedu.io/a"}, URL{"https", "bfedu.io"}, 1},
		{URL{"https", "abc.org"}, URL{"https", "bfedu.io"}, -1},
	}

	for i, tt := range tests {
		result := tt.urlA.Cmp(tt.urlB)
		if result != tt.expect {
			t.Errorf("test %d: cmp mismatch: expected: %d, got: %d", i, tt.expect, result)
		}
	}
}
