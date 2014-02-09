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
	bbTag bbTag
	htmlTag *htmlTag
}

%type <str> tag_end
%type <value> arg
%type <bbTag> tag_start
%type <argument> args
%type <htmlTag> expr
%token <str> TEXT ID NEWLINE EOF

%%

list:
	| list expr
	{ writeExpression(yylex, $2.string()) }
	| list tag_end
	{ writeExpression(yylex, "[/" + $2 + "]") }
	| list EOF
	{ writeExpression(yylex, "") }
	;

expr: tag_start expr tag_end
	{ $$ = compile($1, $2) }
	| expr NEWLINE
	{ $$ = newline($1) }
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

tag_end: '[' '/' ID ']'
	{ $$ = $3 }
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
