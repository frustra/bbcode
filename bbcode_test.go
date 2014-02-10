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
	`[img=http://example.com]alt text[/img]`: `<img src="http://example.com" alt="alt text">`,
	`[img=http://example.com][/img]`:         `<img src="http://example.com">`,

	`[img = foo]bar[/img]`: `<img src="foo" alt="bar">`,

	`[B]bold[/b]`:                    `<b>bold</b>`,
	`[i]italic[/i]`:                  `<i>italic</i>`,
	`[u]underline[/U]`:               `<u>underline</u>`,
	`[strike]strikethrough[/strike]`: `<s>strikethrough</s>`,

	`[u][b]something[/b] then [b]something else[/b][/u]`: `<u><b>something</b> then <b>something else</b></u>`,
	`blank[b][/b]`:                                       `blank<b></b>`,

	"test\nnewline\nnewline": `test<br>newline<br>newline`,
	"test\n\nnewline":        `test<br><br>newline`,
	"[b]test[/b]\n\nnewline": `<b>test</b><br><br>newline`,
	"[b]test\nnewline[/b]":   `<b>test<br>newline</b>`,

	"[code][b]some[/b][i]stuff[/i][/quote][/code][b]more[/b]":                       `<code>[b]some[/b][i]stuff[/i][/quote]</code><b>more</b>`,
	"[center]hello[/center]":                                                        `<div style="text-align: center;">hello</div>`,
	"[size=6]hello[/size]":                                                          `<span style="font-size: 24px;">hello</span>`,
	"[center][b][color=#00BFFF][size=6]hello[/size][/color][/b][/center]":           `<div style="text-align: center;"><b><span style="font-size: 24px;">hello</span></b></div>`,
	"[spoiler][img]http://example.com[/img][/spoiler]":                              `<div class="spoiler-tag"><img src="http://example.com"></div>`,
	"[media]https://www.youtube.com/watch?v=MEQMkzjcLEA&list=RDMEQMkzjcLEA[/media]": `<object width="620" height="349"><param name="wmode" value="transparent"><param name="allowFullScreen" value="true"><param name="allowscriptaccess" value="always"><param name="movie" value="//www.youtube.com/v/MEQMkzjcLEA?version=3"><embed type="application/x-shockwave-flash" width="620" height="349" wmode="transparent" allowFullScreen="true" allowscriptaccess="always" src="//www.youtube.com/v/MEQMkzjcLEA?version=3"></object>`,

	`[not a tag][/not]`: `[not a tag][/not]`,
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
	`[img=<script>]<script>[/img]`: `<img src="%3Cscript%3E" alt="&lt;script&gt;">`,

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

var fullTestInput = `the quick brown [b]fox[/b]:
[url=http://example][img]http://example.png[/img][/url]`

var fullTestOutput = `the quick brown <b>fox</b>:<br><a href="http://example"><img src="http://example.png"></a>`

func TestFull(t *testing.T) {
	result, err := Compile(fullTestInput)
	if err != nil {
		t.Errorf("Unexpected error %v while compiling %s\n", err, fullTestInput)
	} else if result != fullTestOutput {
		t.Errorf("Failed to compile %s.\nExpected: %s, got: %s\n", fullTestInput, fullTestOutput, result)
	}
}

var brokenTests = map[string]string{
	"[b]":        `[b]`,
	"[b]\n":      `[b]<br>`,
	"[b]hello":   `[b]hello`,
	"[b]hello\n": `[b]hello<br>`,
	"the quick brown [b][i]fox[/b][/i]\n[i]\n[b]hi[/b]][b][url=http://example[img]http://example.png[/img][/url][b]": `the quick brown [b][i]fox[/b][/i]<br>[i]<br><b>hi</b>][b][url=http://example[img]http://example.png[/img][/url][b]`,
	"the quick brown[/b][b]hello[/b]":                                                                                `the quick brown[/b]<b>hello</b>`,
	"the quick brown[/b][/code]":                                                                                     `the quick brown[/b][/code]`,
	"[ b][	i]the quick brown[/i][/b=hello]": `[b]<i>the quick brown</i>[/b=hello]`,
	"[b [herp@#$%]]the quick brown[/b]": `<b>]the quick brown</b>`,
	"[b=hello a=hi	q]the quick brown[/b]": `<b>the quick brown</b>`,
	"[b]hi[": `[b]hi[`,
}

func TestBroken(t *testing.T) {
	for in, out := range brokenTests {
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
