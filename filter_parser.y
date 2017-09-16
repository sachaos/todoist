%{
package main

import (
    "text/scanner"
    "strings"
    "strconv"
)

type Expression interface{}
type Token struct {
    token int
    literal string
}

type VoidExpr struct {}

type StringExpr struct {
    literal string
}

type SpecificDateTimeExpr struct {
    year   int
    month  int
    day    int
    hour   int
    minute int
    second int
}

type BoolInfixOpExpr struct {
    left Expression
    operator rune
    right Expression
}

func atoi(a string) (i int) {
    i, _ = strconv.Atoi(a)
    return
}

%}

%union{
    token Token
    expr Expression
}

%type<expr> filter
%type<expr> expr
%type<expr> s_date
%type<expr> s_time
%token<token> STRING
%token<token> NUMBER
%token<token> MONTH_IDENT
%token<token> TWELVE_CLOCK_IDENT
%left '&' '|'

%%

filter
    :
    {
        $$ = VoidExpr{}
    }
    | expr
    {
        $$ = $1
        yylex.(*Lexer).result = $$
    }

expr
    : expr '|' expr
    {
        $$ = BoolInfixOpExpr{left: $1, operator: '|', right: $3}
    }
    | expr '&' expr
    {
        $$ = BoolInfixOpExpr{left: $1, operator: '&', right: $3}
    }
    | STRING
    {
        $$ = StringExpr{literal: $1.literal}
    }
    | '(' expr ')'
    {
        $$ = $2
    }
    | s_date s_time
    {
        date := $1.(SpecificDateTimeExpr)
        time := $2.(SpecificDateTimeExpr)
        $$ = SpecificDateTimeExpr{
            year: date.year, month: date.month, day: date.day,
            hour: time.hour, minute: time.minute, second: time.second,
        }
    }
    | s_date
    {
        $$ = $1
    }
    | s_time
    {
        $$ = $1
    }

s_date
    : NUMBER '/' NUMBER '/' NUMBER
    {
        $$ = SpecificDateTimeExpr{year: atoi($5.literal), month: atoi($1.literal), day: atoi($3.literal)}
    }
    | MONTH_IDENT NUMBER
    {
        $$ = SpecificDateTimeExpr{month: MonthIdentHash[$1.literal], day: atoi($2.literal)}
    }
    | NUMBER MONTH_IDENT
    {
        $$ = SpecificDateTimeExpr{month: MonthIdentHash[$2.literal], day: atoi($1.literal)}
    }
    | MONTH_IDENT NUMBER NUMBER
    {
        $$ = SpecificDateTimeExpr{year: atoi($3.literal), month: MonthIdentHash[$1.literal], day: atoi($2.literal)}
    }
    | NUMBER MONTH_IDENT NUMBER
    {
        $$ = SpecificDateTimeExpr{year: atoi($3.literal), month: MonthIdentHash[$2.literal], day: atoi($1.literal)}
    }
    | NUMBER '/' NUMBER
    {
        $$ = SpecificDateTimeExpr{month: atoi($3.literal), day: atoi($1.literal)}
    }

s_time
    : NUMBER ':' NUMBER
    {
        $$ = SpecificDateTimeExpr{hour: atoi($1.literal), minute: atoi($3.literal)}
    }
    | NUMBER ':' NUMBER ':' NUMBER
    {
        $$ = SpecificDateTimeExpr{hour: atoi($1.literal), minute: atoi($3.literal), second: atoi($5.literal)}
    }
    | NUMBER TWELVE_CLOCK_IDENT
    {
        hour := atoi($1.literal)
        if TwelveClockIdentHash[$2.literal] {
            hour = hour + 12
        }
        $$ = SpecificDateTimeExpr{hour: hour}
    }

%%

type Lexer struct {
    scanner.Scanner
    result Expression
}

var MonthIdentHash = map[string]int{
    "Jan": 1,
    "Feb": 2,
    "Mar": 3,
    "Apr": 4,
    "May": 5,
    "June": 6,
    "July": 7,
    "Aug": 8,
    "Sept": 9,
    "Oct": 10,
    "Nov": 11,
    "Dec": 12,

    "January": 1,
    "February": 2,
    "March": 3,
    "April": 4,
    "August": 8,
    "September": 9,
    "October": 10,
    "November": 11,
    "December": 12,
}

var TwelveClockIdentHash = map[string]bool{
    "am": false,
    "pm": true,
}

func (l *Lexer) Lex(lval *yySymType) int {
    token := int(l.Scan())
    switch token {
        case scanner.Ident:
            if _, ok := MonthIdentHash[l.TokenText()]; ok {
                token = MONTH_IDENT
            } else if _, ok := TwelveClockIdentHash[l.TokenText()]; ok {
                token = TWELVE_CLOCK_IDENT
            } else {
                token = STRING
            }
        case scanner.Int:
            token = NUMBER
    }
    lval.token = Token{token: token, literal: l.TokenText()}
    return token
}

func (l *Lexer) Error(e string) {
    panic(e)
}

func Filter(f string) (e Expression) {
    l := new(Lexer)
    l.Init(strings.NewReader(f))
    yyParse(l)
    return l.result
}
