// Copyright 2015 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import "testing"

var urlTests = map[string]string{
	"http://example.com/path?query=value#fragment":         "http://example.com/path?query=value#fragment",
	"<script>http://example.com":                           "",
	"http://example.com/path?query=value#fragment<script>": "http://example.com/path?query=value#fragment%3Cscript%3E",
	"http://example.com/path?query=<script>":               "http://example.com/path?query=<script>",
	"javascript:alert(1);":                                 "javascript:alert(1);",
}

func TestValidURL(t *testing.T) {
	for in, out := range urlTests {
		result := ValidURL(in)
		if result != out {
			t.Errorf("Failed to sanitize %s.\nExpected: %s, got: %s", in, out, result)
		}
	}
}
