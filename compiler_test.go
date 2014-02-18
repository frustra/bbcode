// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import "testing"

var basicTests = map[string]string{
	``: ``,
	`[url]http://example.com[/url]`: `<a href="http://example.com">http://example.com</a>`,
	`[img]http://example.com[/img]`: `<img src="http://example.com">`,
	`[img][/img]`:                   `<img src="">`,

	`[url=http://example.com]example[/url]`:  `<a href="http://example.com">example</a>`,
	`[img=http://example.com]alt text[/img]`: `<img src="http://example.com" alt="alt text" title="alt text">`,
	`[img=http://example.com][/img]`:         `<img src="http://example.com">`,

	`[img = foo]bar[/img]`: `<img src="foo" alt="bar" title="bar">`,

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

	"[code][b]some[/b][i]stuff[/i][/quote][/code][b]more[/b]":             `<code>[b]some[/b][i]stuff[/i][/quote]</code><b>more</b>`,
	"[quote name=Someguy]hello[/quote]":                                   `<blockquote><cite>Someguy said:</cite>hello</blockquote>`,
	"[center]hello[/center]":                                              `<div style="text-align: center;">hello</div>`,
	"[size=6]hello[/size]":                                                `<span style="font-size: 24px;">hello</span>`,
	"[center][b][color=#00BFFF][size=6]hello[/size][/color][/b][/center]": `<div style="text-align: center;"><b><span style="color: #00BFFF;"><span style="font-size: 24px;">hello</span></span></b></div>`,

	`[not a tag][/not ]`: `[not a tag][/not ]`,
	`[not a tag]`:        `[not a tag]`,
}

func TestCompile(t *testing.T) {
	c := NewCompiler(false, false)
	for in, out := range basicTests {
		result := c.Compile(in)
		if result != out {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
	}
}

var sanitizationTests = map[string]string{
	`<script>`:            `&lt;script&gt;`,
	`[url]<script>[/url]`: `<a href="%3Cscript%3E">&lt;script&gt;</a>`,

	`[url=<script>]<script>[/url]`: `<a href="%3Cscript%3E">&lt;script&gt;</a>`,
	`[img=<script>]<script>[/img]`: `<img src="%3Cscript%3E" alt="&lt;script&gt;" title="&lt;script&gt;">`,

	`[url=http://a.b/z?\]link[/url]`:      `<a href="http://a.b/z?\">link</a>`,
	`[img="http://\"a.b/z"]"link"\[/img]`: `<img src="http://&#34;a.b/z" alt="&#34;link&#34;\" title="&#34;link&#34;\">`,
}

func TestSanitization(t *testing.T) {
	c := NewCompiler(false, false)
	for in, out := range sanitizationTests {
		result := c.Compile(in)
		if result != out {
			t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", in, out, result)
		}
	}
}

var fullTestInput = `the quick brown [b]fox[/b]:
[url=http://example][img]http://example.png[/img][/url]`

var fullTestOutput = `the quick brown <b>fox</b>:<br><a href="http://example"><img src="http://example.png"></a>`

func TestFull(t *testing.T) {
	c := NewCompiler(false, false)
	result := c.Compile(fullTestInput)
	if result != fullTestOutput {
		t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", fullTestInput, fullTestOutput, result)
	}
}

func BenchmarkFull(b *testing.B) {
	c := NewCompiler(false, false)
	for i := 0; i < b.N; i++ {
		c.Compile(fullTestInput)
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

var customTests = map[string]string{
	`[img]//foo/bar.png[/img]`: `<img src="//custom.png">`,
	`[url]//foo/bar.png[/url]`: `[url]//foo/bar.png[/url]`,
}

func compileImg(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
	out, appendExpr := DefaultTagCompilers["img"](node, in)
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
