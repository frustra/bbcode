// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package bbcode implements a parser and HTML generator for BBCode.
package bbcode

// Compile transforms a string of BBCode to HTML
func Compile(str string) string {
	lex := newLexer(str)
	lex.Run()
	return lex.buffer.String()
}

type stringPair struct {
	key, value string
}

type argument struct {
	others *argument
	arg    stringPair
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

func writeExpression(lex yyLexer, expr string) {
	lex.(*lexer).buffer.WriteString(expr)
}

func (a *argument) reduceArguments(args map[string]string) {
	args[a.arg.key] = a.arg.value
	if a.others != nil {
		a.others.reduceArguments(args)
	}
}

func (a *argument) expand() map[string]string {
	var args = make(map[string]string)
	a.reduceArguments(args)
	return args
}
