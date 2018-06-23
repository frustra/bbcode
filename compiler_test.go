// Copyright 2015 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"strings"
	"testing"
)

var fullTestInput = `the quick brown [b]fox[/b]:
[url=http://example][img]http://example.png[/img][/url]`

var fullTestOutput = `the quick brown <b>fox</b>:<br><a href="http://example"><img src="http://example.png"></a>`

func TestFullBasic(t *testing.T) {
	c := NewCompiler(false, false)
	input := fullTestInput
	output := fullTestOutput
	for in, out := range basicTests {
		input += in
		output += out
	}
	result := c.Compile(input)
	if result != output {
		t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", input, output, result)
	}
}

func BenchmarkFullBasic(b *testing.B) {
	c := NewCompiler(false, false)
	input := fullTestInput
	for in := range basicTests {
		input += in
	}
	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Compile(input)
	}
}

var basicTests = map[string]string{
	``: ``,
	`[url]http://example.com[/url]`: `<a href="http://example.com">http://example.com</a>`,
	`[img]http://example.com[/img]`: `<img src="http://example.com">`,
	`[img][/img]`:                   `<img src="">`,

	`[url=http://example.com]example[/url]`: `<a href="http://example.com">example</a>`,
	`[img=http://example.com][/img]`:        `<img src="http://example.com">`,

	`[B]bold[/b]`:          `<b>bold</b>`,
	`[i]italic[/i]`:        `<i>italic</i>`,
	`[u]underline[/U]`:     `<u>underline</u>`,
	`[s]strikethrough[/s]`: `<s>strikethrough</s>`,

	`[u][b]something[/b] then [b]something else[/b][/u]`: `<u><b>something</b> then <b>something else</b></u>`,
	`blank[b][/b]`:                                       `blank<b></b>`,

	"test\nnewline\nnewline": `test<br>newline<br>newline`,
	"test\n\nnewline":        `test<br><br>newline`,
	"[b]test[/b]\n\nnewline": `<b>test</b><br><br>newline`,
	"[b]test\nnewline[/b]":   `<b>test<br>newline</b>`,

	"[code][b]some[/b]\n[i]stuff[/i]\n[/quote][/code][b]more[/b]":         "<pre>[b]some[/b]\n[i]stuff[/i]\n[/quote]</pre><b>more</b>",
	"[quote name=Someguy]hello[/quote]":                                   `<blockquote><cite>Someguy said:</cite>hello</blockquote>`,
	"[center]hello[/center]":                                              `<div style="text-align: center;">hello</div>`,
	"[size=6]hello[/size]":                                                `<span class="size6">hello</span>`,
	"[center][b][color=#00BFFF][size=6]hello[/size][/color][/b][/center]": `<div style="text-align: center;"><b><span style="color: #00BFFF;"><span class="size6">hello</span></span></b></div>`,

	`[not a tag][/not ]`: `[not a tag][/not ]`,
	`[not a tag]`:        `[not a tag]`,
}
var basicMultiArgTests = map[string][]string{
	`[img=http://example.com]alt text[/img]`: []string{`<img`, ` alt="alt text"`, ` src="http://example.com"`, ` title="alt text"`, `>`},
	`[img = foo]bar[/img]`:                   []string{`<img`, ` alt="bar"`, ` src="foo"`, ` title="bar"`, `>`},
}

func TestCompile(t *testing.T) {
	c := NewCompiler(false, false)
	for in, out := range basicTests {
		result := c.Compile(in)
		if result != out {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
	}
	for in, out := range basicMultiArgTests {
		result := c.Compile(in)
		if !strings.HasPrefix(result, out[0]) || !strings.HasSuffix(result, out[len(out)-1]) {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
		for i := 1; i < len(out)-1; i++ {
			if !strings.Contains(result, out[i]) {
				t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
			}
		}
	}
}

func TestSortedAttributes(t *testing.T) {
	c := NewCompiler(false, false)
	c.SortOutputAttributes = true
	// Test this 10 times to eliminate randomness
	for i := 0; i < 10; i++ {
		for in, out := range basicMultiArgTests {
			result := c.Compile(in)
			compare := strings.Join(out, "")
			if result != compare {
				t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, compare, result)
				return
			}
		}
	}
}

var sanitizationTests = map[string]string{
	`<script>`:            `&lt;script&gt;`,
	`[url]<script>[/url]`: `<a href="%3Cscript%3E">&lt;script&gt;</a>`,

	`[url=<script>]<script>[/url]`: `<a href="%3Cscript%3E">&lt;script&gt;</a>`,

	`[url=http://a.b/z?\]link[/url]`: `<a href="http://a.b/z?\">link</a>`,
}
var sanitizationMultiArgTests = map[string][]string{
	`[img=<script>]<script>[/img]`:        []string{`<img`, ` src="%3Cscript%3E"`, ` alt="&lt;script&gt;"`, ` title="&lt;script&gt;"`, `>`},
	`[img="http://\"a.b/z"]"link"\[/img]`: []string{`<img`, ` src="http://&#34;a.b/z"`, ` alt="&#34;link&#34;\"`, ` title="&#34;link&#34;\"`, `>`},
}

func TestSanitization(t *testing.T) {
	c := NewCompiler(false, false)
	for in, out := range sanitizationTests {
		result := c.Compile(in)
		if result != out {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
	}
	for in, out := range sanitizationMultiArgTests {
		result := c.Compile(in)
		if !strings.HasPrefix(result, out[0]) || !strings.HasSuffix(result, out[len(out)-1]) {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
		for i := 1; i < len(out)-1; i++ {
			if !strings.Contains(result, out[i]) {
				t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
			}
		}
	}
}

func TestFullSanitization(t *testing.T) {
	c := NewCompiler(false, false)
	input := fullTestInput
	output := fullTestOutput
	for in, out := range sanitizationTests {
		input += in
		output += out
	}
	result := c.Compile(input)
	if result != output {
		t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", input, output, result)
	}
}

func BenchmarkFullSanitization(b *testing.B) {
	c := NewCompiler(false, false)
	input := fullTestInput
	for in := range sanitizationTests {
		input += in
	}
	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Compile(input)
	}
}

var brokenTests = map[string]string{
	"[b]":        `[b]`,
	"[b]\n":      `[b]<br>`,
	"[b]hello":   `[b]hello`,
	"[b]hello\n": `[b]hello<br>`,
	"the quick brown [b][i]fox[/b][/i]\n[i]\n[b]hi[/b]][b][url=http://example[img]http://example.png[/img][/url][b]": `the quick brown <b>[i]fox</b>[/i]<br>[i]<br><b>hi</b>][b][url=http://example<img src="http://example.png">[/url][b]`,
	"the quick brown[/b][b]hello[/b]":                                                                                `the quick brown[/b]<b>hello</b>`,
	"the quick brown[/b][/code]":                                                                                     `the quick brown[/b][/code]`,
	"[ b][	i]the quick brown[/i][/b=hello]": `[ b]<i>the quick brown</i>[/b=hello]`,
	"[b [herp@#$%]]the quick brown[/b]": `[b [herp@#$%]]the quick brown[/b]`,
	"[b=hello a=hi	q]the quick brown[/b]": `<b>the quick brown</b>`,
	"[b]hi[":     `[b]hi[`,
	"[b hi=derp": `[b hi=derp`,
}

func TestBroken(t *testing.T) {
	c := NewCompiler(false, false)
	for in, out := range brokenTests {
		result := c.Compile(in)
		if result != out {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
	}
}

func BenchmarkFullBroken(b *testing.B) {
	c := NewCompiler(false, false)
	input := fullTestInput
	for in := range brokenTests {
		input += in
	}
	b.SetBytes(int64(len(input)))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Compile(input)
	}
}

var customTests = map[string]string{
	`[img]//foo/bar.png[/img]`: `<img src="//custom.png">`,
	`[url]//foo/bar.png[/url]`: `[url]//foo/bar.png[/url]`,
}

func compileImg(node *BBCodeNode) (*HTMLTag, bool) {
	out, appendExpr := DefaultTagCompilers["img"](node)
	out.Attrs["src"] = "//custom.png"
	return out, appendExpr
}

func TestCompileCustom(t *testing.T) {
	c := NewCompiler(false, false)
	c.SetTag("url", nil)
	c.SetTag("img", compileImg)
	for in, out := range customTests {
		result := c.Compile(in)
		if result != out {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
	}
}
