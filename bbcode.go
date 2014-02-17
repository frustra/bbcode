// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package bbcode implements a parser and HTML generator for BBCode.
package bbcode

// Compile transforms a string of BBCode to HTML.
func Compile(str string) string {
	return CompileCustom(str, DefaultCompiler{})
}

// CompileCustom uses a custom Compiler implementation to transform a string of
// BBCode to HTML.
func CompileCustom(str string, compiler Compiler) string {
	tokens := Lex(str)
	tree := Parse(tokens)
	return compiler.Compile(tree).String()
}

type bbOpeningTag struct {
	Name  string
	Value string
	Args  map[string]string
	Raw   string
}

type bbClosingTag struct {
	Name string
	Raw  string
}

func (t *bbOpeningTag) String() string {
	str := t.Name
	if len(t.Value) > 0 {
		str += "=" + t.Value
	}
	for k, v := range t.Args {
		str += " " + k
		if len(v) > 0 {
			str += "=" + v
		}
	}
	return str
}
