// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"fmt"
	"testing"
)

var prelexTests = map[string][]string{
	``:                                                                    []string{},
	`[url]a[/url]`:                                                        []string{`[url]`, `a`, `[/url]`},
	`[img][/img]`:                                                         []string{`[img]`, `[/img]`},
	`[img = foo]bar[/img]`:                                                []string{`[img=foo]`, `bar`, `[/img]`},
	`[quote name=Someguy]hello[/quote]`:                                   []string{`[quote name=Someguy]`, `hello`, `[/quote]`},
	`[center][b][color=#00BFFF][size=6]hello[/size][/color][/b][/center]`: []string{`[center]`, `[b]`, `[color=#00BFFF]`, `[size=6]`, `hello`, `[/size]`, `[/color]`, `[/b]`, `[/center]`},
	`[b]`:               []string{`[b]`},
	`blank[b][/b]`:      []string{`blank`, `[b]`, `[/b]`},
	`[b][/b]blank`:      []string{`[b]`, `[/b]`, `blank`},
	`[not a tag][/not]`: []string{`[not a tag]`, `[/not]`},

	`[u][b]something[/b] then [b]something else[/b][/u]`: []string{`[u]`, `[b]`, `something`, `[/b]`, ` then `, `[b]`, `something else`, `[/b]`, `[/u]`},

	"the quick brown [b][i]fox[/b][/i]\n[i]\n[b]hi[/b]][b][url=a[img]v[/img][/url][b]": []string{"the quick brown ", "[b]", "[i]", "fox", "[/b]", "[/i]", "\n", "[i]", "\n", "[b]", "hi", "[/b]", "]", "[b]", "[url=a", "[img]", "v", "[/img]", "[/url]", "[b]"},
	"the quick brown[/b][b]hello[/b]":                                                  []string{"the quick brown", "[/b]", "[b]", "hello", "[/b]"},
	"the quick brown[/b][/code]":                                                       []string{"the quick brown", "[/b]", "[/code]"},

	`[ b][	i]the quick brown[/i][/b=hello]`: []string{`[b]`, `[i]`, `the quick brown`, `[/i]`, `[/b=hello]`},
	`[b [herp@#$%]]the quick brown[/b]`: []string{`[b `, `[herp@#$%]`, `]the quick brown`, `[/b]`},
	`[b=hello a=hi	q]the quick brown[/b]`: []string{`[b=hello a=hi q]`, `the quick brown`, `[/b]`},
	`[b]hi[`: []string{`[b]`, `hi[`},

	`[img = 'fo"o']bar[/img]`:                                                []string{`[img=fo"o]`, `bar`, `[/img]`},
	`[img = "foo'"]bar[/img]`:                                                []string{`[img=foo']`, `bar`, `[/img]`},
	`[img = "\"'foo"]bar[/img]`:                                              []string{`[img=\ 'foo"]`, `bar`, `[/img]`},
	`[quote name='Someguy']hello[/quote]`:                                    []string{`[quote name=Someguy]`, `hello`, `[/quote]`},
	`[center][b][color="#00BFFF"][size='6]hello[/size][/color][/b][/center]`: []string{`[center]`, `[b]`, `[color=#00BFFF]`, `[size=6]`, `hello`, `[/size]`, `[/color]`, `[/b]`, `[/center]`},
}

func TestPreLex(t *testing.T) {
	for in, expected := range prelexTests {
		lexer := newLexer(in)
		lexer.PreLex()
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

func PrintOutput(out []*token) string {
	result := ""
	for i, v := range out {
		if i > 0 {
			result += "_"
		}
		switch t := v.value.(type) {
		case string:
			result += t
		case bbOpeningTag:
			result += t.string()
		case bbClosingTag:
			result += "[/" + t.name + "]"
		default:
			result += fmt.Sprintf("<%v>", t)
		}
	}
	return result
}

func CheckResult(l *lexer, b []string) (bool, []*token) {
	i := 0
	out := make([]*token, 0)
	good := true
	for v := range l.tokens {
		if v == nil {
			break
		}
		out = append(out, v)
		if i < len(b) && good {
			switch t := v.value.(type) {
			case string:
				if t != b[i] {
					good = false
				}
			case bbOpeningTag:
				if t.string() != b[i] {
					good = false
				}
			case bbClosingTag:
				if "[/"+t.name+"]" != b[i] {
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
