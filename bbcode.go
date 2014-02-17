// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package bbcode implements a parser and HTML generator for BBCode.
package bbcode

// Compile transforms a string of BBCode to HTML
func Compile(str string) string {
	tokens := Lex(str)
	tree := Parse(tokens)
	return compile(tree).string()
}

type bbOpeningTag struct {
	name  string
	value string
	args  map[string]string
	raw   string
}

type bbClosingTag struct {
	name string
	raw  string
}

func (t *bbOpeningTag) string() string {
	str := t.name
	if len(t.value) > 0 {
		str += "=" + t.value
	}
	for k, v := range t.args {
		str += " " + k
		if len(v) > 0 {
			str += "=" + v
		}
	}
	return str
}
