// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import "testing"

var urlTests = map[string]string{
	"http://example.com/path?query=value#fragment": "http://example.com/path?query=value#fragment",
	"<script>http://example.com":                   "%3Cscript%3Ehttp://example.com",
}

func TestSafeURL(t *testing.T) {
	for in, out := range urlTests {
		result := safeURL(in)
		if result != out {
			t.Errorf("Failed to sanitize %s.\nExpected: %s, got: %s", in, out, result)
		}
	}
}
