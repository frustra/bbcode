// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

%{
package bbcode
%}

%union{
	str string
	value stringPair
	argument *argument
	openingTag bbOpeningTag
	closingTag bbClosingTag
	htmlTag *htmlTag
}

%type <htmlTag> list
%type <htmlTag> expr
%token <str> TEXT MISSING_CLOSING MISSING_OPENING ID NEWLINE
%token <openingTag> OPENING
%token <closingTag> CLOSING

%%

full: list
	{
		if $1 != nil {
			writeExpression(yylex, $1.string())
		}
	}
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

expr: OPENING list CLOSING
	{
		compile($1, $2)
	}
	| OPENING list MISSING_CLOSING
	{ $$ = newHtmlTag($1.string()).appendChild($2) }
	| MISSING_OPENING ID ']'
	{ $$ = newHtmlTag("[/" + $2 + "]") }
	| NEWLINE
	{ $$ = newline() }
	| TEXT
	{ $$ = newHtmlTag($1) }
	;

%%
