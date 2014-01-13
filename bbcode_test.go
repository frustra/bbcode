// Copyright 2014 Frustra Sofware. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"testing"
)

var basicTests = map[string]string{
	`[url=http://example.com]example[/url]`:  `<a href="http://example.com">example</a>`,
	`[url]http://example.com[/url]`:          `<a href="http://example.com">http://example.com</a>`,
	`[img]http://example.com[/img]`:          `<img src="http://example.com"/>`,
	`[img=http://example.com]alt text[/img]`: `<img src="http://example.com" alt="alt text"/>`,
}

func TestCompile(t *testing.T) {
	for in, out := range basicTests {
		result, err := Compile(in)
		if err != nil {
			t.Errorf("Unexpected error %v while compiling %s\n", err, in)
		} else if result != out {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
	}
}

var sanitizationTests = map[string]string{
	`<script>`:            `&lt;script&gt;`,
	`[url]<script>[/url]`: `<a href="%3Cscript%3E">&lt;script&gt;</a>`,

	`[url=<script>]<script>[/url]`: `<a href="%3Cscript%3E">&lt;script&gt;</a>`,
	`[img=<script>]<script>[/url]`: `<img src="%3Cscript%3E" alt="&lt;script&gt;"/>`,

	`[url=http://a.b/z?\]link[/url]`: `<a href="http://a.b/z?%5C">link</a>`,
}

func TestSanitization(t *testing.T) {
	for in, out := range sanitizationTests {
		result, err := Compile(in)
		if err != nil {
			t.Errorf("Unexpected error %v while compiling %s\n", err, in)
		} else if result != out {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
	}
}

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
