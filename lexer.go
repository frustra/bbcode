// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

type lexerState int

const (
	INIT_STATE lexerState = iota
	TAG_START_STATE
	TAG_ARGS_STATE
	ARG_VALUE_STATE
)

// Abuse of the generated types to keep parser state in the lexer
type lexer struct {
	str        []byte
	state      lexerState
	tagsOpened int
	buffer     bytes.Buffer
	err        error
}

var (
	tags            = []string{"url", "img", "b", "i", "u", "s", "quote", "code"}
	tagOpenRegexps  []*regexp.Regexp
	tagCloseRegexps []*regexp.Regexp
)

func init() {
	for _, tag := range tags {
		r := regexp.MustCompile(`(?i)^\[[ \t]*` + tag + `[\]= \t]`)
		tagOpenRegexps = append(tagOpenRegexps, r)
		r = regexp.MustCompile(`(?i)^\[/[ \t]*` + tag + `[ \t]*\]`)
		tagCloseRegexps = append(tagCloseRegexps, r)
	}
}

func newLexer(str string) *lexer {
	return &lexer{
		str: []byte(str),
	}
}

func (l *lexer) Lex(lval *yySymType) int {
	if len(l.str) <= 0 {
		if l.tagsOpened > 0 {
			l.tagsOpened--
			return MISSING_CLOSING
		} else {
			return 0
		}
	}
	var c byte = l.str[0]

	switch l.state {
	case TAG_START_STATE:
		for unicode.IsSpace(rune(c)) {
			l.str = l.str[1:]
			c = l.str[0]
		}
		str := strings.ToLower(string(l.str))
		for _, tag := range tags {
			if strings.HasPrefix(str, tag) {
				lval.str = tag
				l.str = l.str[len(tag):]
				l.state = TAG_ARGS_STATE
				return ID
			}
		}
	case TAG_ARGS_STATE:
		for unicode.IsSpace(rune(c)) {
			l.str = l.str[1:]
			c = l.str[0]
		}
		switch {
		case c == ']':
			l.str = l.str[1:]
			l.state = INIT_STATE
			return int(c)
		case c == '=':
			l.str = l.str[1:]
			l.state = ARG_VALUE_STATE
			return int(c)
		default:
			offset := 1
			for offset < len(l.str) {
				curr := l.str[offset]
				if curr == ']' || curr == '=' {
					break
				}
				offset++
			}
			lval.str = string(l.str[0:offset])
			l.str = l.str[offset:]
			return ID
		}
	case ARG_VALUE_STATE:
		for unicode.IsSpace(rune(c)) {
			l.str = l.str[1:]
			c = l.str[0]
		}
		switch {
		case c == '"' || c == '\'':
			return 0 //l.LexQuotedString(c, lval)
		}
		offset := 1
		for offset < len(l.str) {
			curr := l.str[offset]
			if curr == ']' || curr == ' ' || curr == '\t' {
				break
			}
			offset++
		}
		lval.str = string(l.str[0:offset])
		l.str = l.str[offset:]
		l.state = TAG_ARGS_STATE
		return TEXT
	case INIT_STATE:
		if c == '\n' {
			l.str = l.str[1:]
			return NEWLINE
		}
		if c == '[' {
			if l.str[1] == '/' {
				for _, r := range tagCloseRegexps {
					if r.Match(l.str) {
						l.str = l.str[2:]
						l.state = TAG_START_STATE
						if l.tagsOpened <= 0 {
							return MISSING_OPENING
						} else {
							l.tagsOpened--
							return CLOSING_TAG_OPENING
						}
					}
				}
			} else {
				for _, r := range tagOpenRegexps {
					if r.Match(l.str) {
						l.str = l.str[1:]
						l.state = TAG_START_STATE
						l.tagsOpened++
						return int(c)
					}
				}
			}
		}
		offset := 1
		for offset < len(l.str) {
			curr := l.str[offset]
			if curr == '[' || curr == '\n' {
				break
			}
			offset++
		}
		lval.str = string(l.str[0:offset])
		l.str = l.str[offset:]
		return TEXT
	}
	return TEXT
}

func (l *lexer) Error(s string) {
	l.err = errors.New(s)
}
