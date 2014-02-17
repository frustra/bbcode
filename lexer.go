// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import "bytes"

type stateFn func(*lexer) stateFn

type token struct {
	id    int
	value interface{}
}

// Abuse of the generated types to keep parser state in the lexer
type lexer struct {
	input  string
	tokens chan *token

	start int
	end   int
	pos   int

	tagName     string
	tagValue    string
	tagTmpName  string
	tagTmpValue string
	tagArgs     map[string]string

	buffer bytes.Buffer
}

var tags = []string{"url", "img", "b", "i", "u", "strike", "center", "color", "size", "quote", "code", "spoiler", "media"}

func newLexer(str string) *lexer {
	return &lexer{
		input:  str,
		tokens: make(chan *token),
	}
}

func (l *lexer) emit(id int, value interface{}) {
	if l.pos > 0 {
		// fmt.Println(l.input)
		// fmt.Printf("Emit %s: %+v\n", yyToknames[id-TEXT], value)
		l.tokens <- &token{id, value}
		l.input = l.input[l.pos:]
		l.pos = 0
	}
}

func lexText(l *lexer) stateFn {
	for l.pos < len(l.input) {
		if l.input[l.pos] == '[' {
			l.emit(TEXT, l.input[:l.pos])
			return lexOpenBracket
		}
		l.pos++
	}
	l.emit(TEXT, l.input)
	return nil
}

func lexOpenBracket(l *lexer) stateFn {
	l.pos++
	closingTag := false
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case '[', ']':
			return lexText
		default:
			if l.input[l.pos] == '/' && !closingTag {
				closingTag = true
			} else if l.input[l.pos] != ' ' && l.input[l.pos] != '\t' && l.input[l.pos] != '\n' {
				if closingTag {
					return lexClosingTag
				} else {
					l.tagName = ""
					l.tagValue = ""
					l.tagArgs = make(map[string]string)
					return lexTagName
				}
			}
		}
		l.pos++
	}
	l.emit(TEXT, l.input)
	return nil
}

func lexClosingTag(l *lexer) stateFn {
	whiteSpace := false
	l.start = l.pos
	l.end = l.pos
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case '[':
			return lexText
		case ']':
			l.pos++
			l.emit(CLOSING, bbClosingTag{l.input[l.start:l.end], l.input[:l.pos]})
			return lexText
		case ' ', '\t', '\n':
			whiteSpace = true
		default:
			if whiteSpace {
				return lexText
			} else {
				l.end++
			}
		}
		l.pos++
	}
	l.emit(TEXT, l.input)
	return nil
}

func lexTagName(l *lexer) stateFn {
	l.tagTmpValue = ""
	whiteSpace := false
	l.start = l.pos
	l.end = l.pos
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case '[':
			return lexText
		case ']':
			l.tagTmpName = l.input[l.start:l.end]
			return lexTagArgs
		case '=':
			l.tagTmpName = l.input[l.start:l.end]
			return lexTagValue
		case ' ', '\t', '\n':
			whiteSpace = true
		default:
			if whiteSpace {
				l.tagTmpName = l.input[l.start:l.end]
				return lexTagArgs
			} else {
				l.end++
			}
		}
		l.pos++
	}
	l.emit(TEXT, l.input)
	return nil
}

func lexTagValue(l *lexer) stateFn {
	l.pos++
loop:
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\n':
			l.pos++
		case '"', '\'':
			return lexQuotedValue
		default:
			break loop
		}
	}
	l.start = l.pos
	l.end = l.pos
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case '[':
			return lexText
		case ']':
			l.tagTmpValue = l.input[l.start:l.end]
			return lexTagArgs
		case ' ', '\t', '\n':
			l.tagTmpValue = l.input[l.start:l.end]
			return lexTagArgs
		default:
			l.end++
		}
		l.pos++
	}
	l.emit(TEXT, l.input)
	return nil
}

func lexQuotedValue(l *lexer) stateFn {
	quoteChar := l.input[l.pos]
	l.pos++
	l.start = l.pos
	var buf bytes.Buffer
	escape := false
	for l.pos < len(l.input) {
		if escape {
			if l.input[l.pos] == 'n' {
				buf.WriteRune('\n')
			} else {
				buf.WriteRune(rune(l.input[l.pos]))
			}
			escape = false
		} else {
			switch l.input[l.pos] {
			case '\\':
				escape = true
			case '\n':
				l.pos = l.start
				return lexText
			case quoteChar:
				l.pos++
				l.tagTmpValue = buf.String()
				return lexTagArgs
			default:
				buf.WriteRune(rune(l.input[l.pos]))
			}
		}
		l.pos++
	}
	l.pos = l.start
	return lexText
}

func lexTagArgs(l *lexer) stateFn {
	if len(l.tagName) > 0 {
		l.tagArgs[l.tagTmpName] = l.tagTmpValue
	} else {
		l.tagName = l.tagTmpName
		l.tagValue = l.tagTmpValue
	}
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case '[':
			return lexText
		case ']':
			l.pos++
			l.emit(OPENING, bbOpeningTag{l.tagName, l.tagValue, l.tagArgs, l.input[:l.pos]})
			return lexText
		case ' ', '\t', '\n':
			l.pos++
		default:
			l.tagTmpName = ""
			return lexTagName
		}
	}
	l.emit(TEXT, l.input)
	return nil
}

func (l *lexer) runStateMachine() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) Run() {
	go l.runStateMachine()
	yyParse(l)
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
	panic(s)
}
