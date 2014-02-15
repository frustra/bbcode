// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package bbcode implements a parser and HTML generator for BBCode.
package bbcode

import (
	"fmt"
	"strings"
)

// Compile transforms a string of BBCode to HTML
func Compile(str string) string {
	lex := newLexer(str)
	lex.PreLex()
	yyParse(lex)
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
}

type bbClosingTag struct {
	name string
}

func parseBBCodeTag(tag string) interface{} {
	if len(tag) <= 0 || tag[0] != '[' || tag[len(tag)-1] != ']' {
		fmt.Println("Halp, what's going on? %s", tag)
		return nil
	}
	str := strings.Trim(tag[1:len(tag)-1], " \t")
	if len(str) <= 0 {
		return tag
	}
	if str[0] == '/' {
		return bbClosingTag{str[1:]}
	} else {
		var out bbOpeningTag
		out.args = make(map[string]string)
		nameSet := false
		for len(str) > 0 {
			name := ""
			value := ""
			i := strings.IndexAny(str, "= \t")
			if i < 0 {
				i = len(str)
			}
			name = str[0:i]
			str = strings.TrimLeft(str[i:], " \t")
			if len(str) > 0 && str[0] == '=' {
				str = strings.TrimLeft(str[1:], " \t")
				if len(str) > 0 && str[0] == '\'' || str[0] == '"' {
					i = strings.IndexRune(str[1:], rune(str[0]))
					if i < 0 {
						value = str[1:]
						str = ""
					} else {
						value = str[1 : i+1]
						str = strings.TrimLeft(str[i+2:], " \t")
					}
				} else {
					i = strings.IndexAny(str, " \t")
					if i < 0 {
						value = str
						str = ""
					} else {
						value = str[0:i]
						str = strings.TrimLeft(str[i:], " \t")
					}
				}
			}
			if nameSet {
				out.args[name] = value
			} else {
				out.name = name
				out.value = value
				nameSet = true
			}
		}
		return out
	}
	return tag
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
	return "[" + str + "]"
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
