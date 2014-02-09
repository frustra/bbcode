// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

%{
package bbcode

import "strings"
%}

%union{
	str string
	value stringPair
	argument *argument
	bbTag bbTag
	htmlTag *htmlTag
}

%type <str> tag_end
%type <value> arg
%type <bbTag> tag_start
%type <argument> args
%type <htmlTag> list
%type <htmlTag> expr
%token <str> TEXT ID NEWLINE MISSING_CLOSING CLOSING_TAG_OPENING MISSING_OPENING

%%

full: list
	{ writeExpression(yylex, $1.string()) }
	;

list:
	{ $$ = nil }
	| expr list
	{
		if $2 == nil {
			$$ = $1
		} else {
			$$ = newHtmlTag("").appendChild($1).appendChild($2)
		}
	}
	;

expr: tag_start list tag_end
	{
		if strings.EqualFold($1.key, $3) {
			$$ = compile($1, $2)
		} else {
			$$ = newHtmlTag($1.string()).appendChild($2).appendChild(newHtmlTag("[/" + $3 + "]"))
		}
	}
	| tag_start list MISSING_CLOSING
	{ $$ = newHtmlTag($1.string()).appendChild($2) }
	| MISSING_OPENING ID ']'
	{ $$ = newHtmlTag("[/" + $2 + "]") }
	| NEWLINE
	{ $$ = newline() }
	| TEXT
	{ $$ = newHtmlTag($1) }
	;

tag_start: '[' arg args ']'
	{
		$$.key = $2.key
		$$.value = $2.value
		if $3 != nil {
			$$.args = $3.expand()
		}
	}
	;

tag_end: CLOSING_TAG_OPENING ID ']'
	{ $$ = $2 }
	;

arg: ID
	{ $$.key = $1 }
	| ID '=' TEXT
	{
		$$.key = $1
		$$.value = $3
	}
	;

args:
	{ $$ = nil }
	| args arg
	{
		$$ = &argument{}
		$$.others = $1
		$$.arg = $2
	}
	;

%%
