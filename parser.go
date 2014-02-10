// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// This file is generated from parser.y

//line parser.y:6
package bbcode

import __yyfmt__ "fmt"

//line parser.y:6
import "strings"

//line parser.y:11
type yySymType struct {
	yys      int
	str      string
	value    stringPair
	argument *argument
	bbTag    bbTag
	htmlTag  *htmlTag
}

const TEXT = 57346
const ID = 57347
const NEWLINE = 57348
const MISSING_CLOSING = 57349
const CLOSING_TAG_OPENING = 57350
const MISSING_OPENING = 57351

var yyToknames = []string{
	"TEXT",
	"ID",
	"NEWLINE",
	"MISSING_CLOSING",
	"CLOSING_TAG_OPENING",
	"MISSING_OPENING",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line parser.y:96

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 15
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 25

var yyAct = []int{

	18, 7, 20, 6, 25, 22, 5, 19, 8, 17,
	15, 16, 13, 2, 21, 11, 12, 9, 10, 24,
	23, 1, 3, 4, 14,
}
var yyPact = []int{

	-3, -1000, -1000, -3, -3, 10, -1000, -1000, 7, -1000,
	3, -1, 7, -10, -1000, -1000, 9, -1000, -5, 7,
	15, -6, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 24, 7, 23, 0, 13, 22, 21,
}
var yyR1 = []int{

	0, 7, 5, 5, 6, 6, 6, 6, 6, 3,
	1, 2, 2, 4, 4,
}
var yyR2 = []int{

	0, 1, 0, 2, 3, 3, 3, 1, 1, 4,
	3, 1, 3, 0, 2,
}
var yyChk = []int{

	-1000, -7, -5, -6, -3, 9, 6, 4, 11, -5,
	-5, 5, -2, 5, -1, 7, 8, 10, -4, -2,
	12, 5, 10, -4, 4, 10,
}
var yyDef = []int{

	2, -2, 1, 2, 2, 0, 7, 8, 0, 3,
	0, 0, 13, 11, 4, 5, 0, 6, 0, 13,
	0, 0, 9, 14, 12, 10,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 12, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 11, 3, 10,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line parser.y:30
		{
			writeExpression(yylex, yyS[yypt-0].htmlTag.string())
		}
	case 2:
		//line parser.y:34
		{
			yyVAL.htmlTag = nil
		}
	case 3:
		//line parser.y:36
		{
			if yyS[yypt-0].htmlTag == nil {
				yyVAL.htmlTag = yyS[yypt-1].htmlTag
			} else {
				yyVAL.htmlTag = newHtmlTag("").appendChild(yyS[yypt-1].htmlTag).appendChild(yyS[yypt-0].htmlTag)
			}
		}
	case 4:
		//line parser.y:46
		{
			if strings.EqualFold(yyS[yypt-2].bbTag.key, yyS[yypt-0].str) {
				yyVAL.htmlTag = compile(yyS[yypt-2].bbTag, yyS[yypt-1].htmlTag)
			} else {
				yyVAL.htmlTag = newHtmlTag(yyS[yypt-2].bbTag.string()).appendChild(yyS[yypt-1].htmlTag).appendChild(newHtmlTag("[/" + yyS[yypt-0].str + "]"))
			}
		}
	case 5:
		//line parser.y:54
		{
			yyVAL.htmlTag = newHtmlTag(yyS[yypt-2].bbTag.string()).appendChild(yyS[yypt-1].htmlTag)
		}
	case 6:
		//line parser.y:56
		{
			yyVAL.htmlTag = newHtmlTag("[/" + yyS[yypt-1].str + "]")
		}
	case 7:
		//line parser.y:58
		{
			yyVAL.htmlTag = newline()
		}
	case 8:
		//line parser.y:60
		{
			yyVAL.htmlTag = newHtmlTag(yyS[yypt-0].str)
		}
	case 9:
		//line parser.y:64
		{
			yyVAL.bbTag.key = yyS[yypt-2].value.key
			yyVAL.bbTag.value = yyS[yypt-2].value.value
			if yyS[yypt-1].argument != nil {
				yyVAL.bbTag.args = yyS[yypt-1].argument.expand()
			}
		}
	case 10:
		//line parser.y:74
		{
			yyVAL.str = yyS[yypt-1].str
		}
	case 11:
		//line parser.y:78
		{
			yyVAL.value.key = yyS[yypt-0].str
		}
	case 12:
		//line parser.y:80
		{
			yyVAL.value.key = yyS[yypt-2].str
			yyVAL.value.value = yyS[yypt-0].str
		}
	case 13:
		//line parser.y:87
		{
			yyVAL.argument = nil
		}
	case 14:
		//line parser.y:89
		{
			yyVAL.argument = &argument{}
			yyVAL.argument.others = yyS[yypt-0].argument
			yyVAL.argument.arg = yyS[yypt-1].value
		}
	}
	goto yystack /* stack new state and value */
}
