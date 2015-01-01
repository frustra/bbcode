// Copyright 2015 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"fmt"
	"testing"
)

var prelexTests = map[string][]string{
	``:                                                                    []string{},
	`[url]a[/url]`:                                                        []string{`<url>`, `a`, `</url>`},
	`[img][/img]`:                                                         []string{`<img>`, `</img>`},
	`[img = foo]bar[/img]`:                                                []string{`<img=foo>`, `bar`, `</img>`},
	`[quote name=Someguy]hello[/quote]`:                                   []string{`<quote name=Someguy>`, `hello`, `</quote>`},
	`[center][b][color=#00BFFF][size=6]hello[/size][/color][/b][/center]`: []string{`<center>`, `<b>`, `<color=#00BFFF>`, `<size=6>`, `hello`, `</size>`, `</color>`, `</b>`, `</center>`},
	`[b]`:               []string{`<b>`},
	`blank[b][/b]`:      []string{`blank`, `<b>`, `</b>`},
	`[b][/b]blank`:      []string{`<b>`, `</b>`, `blank`},
	`[not a tag][/not]`: []string{`<not a tag>`, `</not>`},

	`[u][b]something[/b] then [b]something else[/b][/u]`: []string{`<u>`, `<b>`, `something`, `</b>`, ` then `, `<b>`, `something else`, `</b>`, `</u>`},

	"the quick brown [b][i]fox[/b][/i]\n[i]\n[b]hi[/b]][b][url=a[img]v[/img][/url][b]": []string{"the quick brown ", "<b>", "<i>", "fox", "</b>", "</i>", "\n", "<i>", "\n", "<b>", "hi", "</b>", "]", "<b>", "[url=a", "<img>", "v", "</img>", "</url>", "<b>"},
	"the quick brown[/b][b]hello[/b]":                                                  []string{"the quick brown", "</b>", "<b>", "hello", "</b>"},
	"the quick brown[/b][/code]":                                                       []string{"the quick brown", "</b>", "</code>"},
	"[quote\n		name=xthexder\n		time=555555\n	]hello[/quote]": []string{`<quote name=xthexder time=555555>`, `hello`, `</quote>`},
	"[q\nuot\ne\nna\nme\n=\nxthex\nder\n]hello[/quote]": []string{`<q der e me=xthex na uot>`, `hello`, `</quote>`},

	`[ b][	i]the quick brown[/i][/b=hello]`: []string{`<b>`, `<i>`, `the quick brown`, `</i>`, `</b=hello>`},
	`[b [herp@#$%]]the quick brown[/b]`: []string{`[b `, `<herp@#$%>`, `]the quick brown`, `</b>`},
	`[b=hello a=hi	q]the quick brown[/b]`: []string{`<b=hello a=hi q>`, `the quick brown`, `</b>`},
	`[b]hi[`:                       []string{`<b>`, `hi`, `[`},
	`[size=6 =hello]hi[/size]`:     []string{`<size=6 =hello>`, `hi`, `</size>`},
	`[size=6 =hello =hi]hi[/size]`: []string{`<size=6 =hi>`, `hi`, `</size>`},

	`[img = 'fo"o']bar[/img]`:                                                    []string{`<img=fo"o>`, `bar`, `</img>`},
	`[img = "foo'"]bar[/img]`:                                                    []string{`<img=foo'>`, `bar`, `</img>`},
	`[img = "\"'foo"]bar[/img]`:                                                  []string{`<img="'foo>`, `bar`, `</img>`},
	`[img = "f\oo\]\'fo\\o"]bar[/img]`:                                           []string{`<img=foo]'fo\o>`, `bar`, `</img>`},
	`[img = "foo\]'fo\n\"o"]bar[/img]`:                                           []string{"<img=foo]'fo\n\"o>", `bar`, `</img>`},
	`[quote name='Someguy']hello[/quote]`:                                        []string{`<quote name=Someguy>`, `hello`, `</quote>`},
	`[center][b][color="#00BFFF"][size='6]hello[/size][/color][/b][/center]`:     []string{`<center>`, `<b>`, `<color=#00BFFF>`, `[size='6]hello`, `</size>`, `</color>`, `</b>`, `</center>`},
	"[center][b][color=\"#00BFFF\"][size='6]hello[/size]\n[/color][/b][/center]": []string{`<center>`, `<b>`, `<color=#00BFFF>`, `[size='6]hello`, `</size>`, "\n", `</color>`, `</b>`, `</center>`},
}

func TestLexer(t *testing.T) {
	for in, expected := range prelexTests {
		lexer := newLexer(in)
		go lexer.runStateMachine()
		ok, out := CheckResult(lexer, expected)
		if !ok {
			t.Errorf("Failed to prelex %s.\nExpected: %s, got: %s\n", in, PrintExpected(expected), PrintOutput(out))
		}
	}
}

func PrintExpected(expected []string) string {
	result := ""
	for i, v := range expected {
		if i > 0 {
			result += "_"
		}
		result += v
	}
	return result
}

func PrintOutput(out []Token) string {
	result := ""
	for i, v := range out {
		if i > 0 {
			result += "_"
		}
		switch t := v.Value.(type) {
		case string:
			result += t
		case BBOpeningTag:
			result += "<" + t.String() + ">"
		case BBClosingTag:
			result += "</" + t.Name + ">"
		default:
			result += fmt.Sprintf("{%v}", t)
		}
	}
	return result
}

func CheckResult(l *lexer, b []string) (bool, []Token) {
	i := 0
	out := make([]Token, 0)
	good := true
	for v := range l.tokens {
		out = append(out, v)
		if i < len(b) && good {
			switch t := v.Value.(type) {
			case string:
				if t != b[i] {
					good = false
				}
			case BBOpeningTag:
				if "<"+t.String()+">" != b[i] {
					good = false
				}
			case BBClosingTag:
				if "</"+t.Name+">" != b[i] {
					good = false
				}
			default:
				good = false
			}
		}
		i++
	}
	if i != len(b) {
		return false, out
	}
	return good, out
}
