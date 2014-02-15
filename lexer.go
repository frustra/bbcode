// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type token struct {
	id    int
	value interface{}
}

// Abuse of the generated types to keep parser state in the lexer
type lexer struct {
	input  string
	tokens chan *token
	buffer bytes.Buffer
}

var (
	tags            = []string{"url", "img", "b", "i", "u", "strike", "center", "color", "size", "quote", "code", "spoiler", "media"}
	tagOpenRegexps  []*regexp.Regexp
	tagCloseRegexps []*regexp.Regexp
	codeCloseRegex  *regexp.Regexp
)

func init() {
	for _, tag := range tags {
		r := regexp.MustCompile(`(?i)^\[[ \t]*` + tag + `[\]= \t]`)
		tagOpenRegexps = append(tagOpenRegexps, r)
		r = regexp.MustCompile(`(?i)^\[/[ \t]*` + tag + `[ \t]*\]`)
		if tag == "code" {
			codeCloseRegex = r
		} else {
			tagCloseRegexps = append(tagCloseRegexps, r)
		}
	}
}

func newLexer(str string) *lexer {
	return &lexer{
		input:  str,
		tokens: make(chan *token, 1000),
	}
}

func (l *lexer) PreLex() {
	for len(l.input) > 0 {
		i := strings.IndexRune(l.input, '[')
		if i < 0 {
			l.tokens <- &token{TEXT, l.input}
			l.input = ""
			continue
		}
		j := strings.IndexAny(l.input[i+1:], "[]")
		if j < 0 {
			l.tokens <- &token{TEXT, l.input}
			l.input = ""
			continue
		}
		if l.input[i+j+1] == ']' {
			if i == 0 {
				parsed := parseBBCodeTag(l.input[0 : j+2])
				switch t := parsed.(type) {
				case nil:
				case string:
					l.tokens <- &token{TEXT, t}
				case bbOpeningTag:
					l.tokens <- &token{OPENING, t}
				case bbClosingTag:
					l.tokens <- &token{CLOSING, t}
				default:
					fmt.Println("Unknown type back from lexer:", t)
				}
				l.input = l.input[j+2:]
			} else {
				l.tokens <- &token{TEXT, l.input[0:i]}
				l.input = l.input[i:]
			}
		} else {
			l.tokens <- &token{TEXT, l.input[0 : j+1]}
			l.input = l.input[j+1:]
		}
	}
	l.tokens <- nil
	close(l.tokens)
}

func (l *lexer) Lex(lval *yySymType) int {
	token := <-l.tokens
	if token != nil {
		switch t := token.value.(type) {
		case string:
			lval.str = t
		case bbOpeningTag:
			lval.openingTag = t
		case bbClosingTag:
			lval.closingTag = t
		}
		return token.id
	} else {
		return 0
	}
}

func (l *lexer) Error(s string) {
	panic(errors.New(s))
}
