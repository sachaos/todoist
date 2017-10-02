%{
package main

import (
    "text/scanner"
    "strings"
    "strconv"
    "time"
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

type BoolInfixOpExpr struct {
    left Expression
    operator rune
    right Expression
}

func atoi(a string) (i int) {
    i, _ = strconv.Atoi(a)
    return
}

var now = time.Now
var today = func() time.Time {
  return time.Date(now().Year(), now().Month(), now().Day(), 0, 0, 0, 0, time.Local)
}

%}

%union{
    token Token
    expr Expression
}

%type<expr> filter
%type<expr> expr
%type<expr> s_datetime
%type<expr> s_date
%type<expr> s_date_year
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
    | s_datetime

s_datetime
    : s_date_year s_time
    {
        date := $1.(time.Time)
        time := $2.(time.Duration)
        $$ = date.Add(time)
    }
    | s_date_year
    {
        $$ = $1
    }
    | s_time
    {
        nd := now().Sub(today())
        d := $1.(time.Duration)
        if (d <= nd) {
          d = d + time.Duration(int64(time.Hour) * 24)
        }
        $$ = today().Add(d)
    }

s_date_year
    : NUMBER '/' NUMBER '/' NUMBER
    {
        $$ = time.Date(atoi($5.literal), time.Month(atoi($1.literal)), atoi($3.literal), 0, 0, 0, 0, time.Local)
    }
    | MONTH_IDENT NUMBER NUMBER
    {
        $$ = time.Date(atoi($3.literal), MonthIdentHash[$1.literal], atoi($2.literal), 0, 0, 0, 0, time.Local)
    }
    | NUMBER MONTH_IDENT NUMBER
    {
        $$ = time.Date(atoi($3.literal), MonthIdentHash[$2.literal], atoi($1.literal), 0, 0, 0, 0, time.Local)
    }
    | s_date
    {
        tod := today()
        date := $1.(time.Time)
        if date.Before(tod) {
            date = date.AddDate(1, 0, 0)
        }
        $$ = date
    }

s_date
    : MONTH_IDENT NUMBER
    {
        $$ = time.Date(today().Year(), MonthIdentHash[$1.literal], atoi($2.literal), 0, 0, 0, 0, time.Local)
    }
    | NUMBER MONTH_IDENT
    {
        $$ = time.Date(today().Year(), MonthIdentHash[$2.literal], atoi($1.literal), 0, 0, 0, 0, time.Local)
    }
    | NUMBER '/' NUMBER
    {
        $$ = time.Date(now().Year(), time.Month(atoi($3.literal)), atoi($1.literal), 0, 0, 0, 0, time.Local)
    }

s_time
    : NUMBER ':' NUMBER
    {
        $$ = time.Duration(int64(time.Hour) * int64(atoi($1.literal)) + int64(time.Minute) * int64(atoi($3.literal)))
    }
    | NUMBER ':' NUMBER ':' NUMBER
    {
        $$ = time.Duration(int64(time.Hour) * int64(atoi($1.literal)) + int64(time.Minute) * int64(atoi($3.literal)) + int64(time.Second) * int64(atoi($5.literal)))
    }
    | NUMBER TWELVE_CLOCK_IDENT
    {
        hour := atoi($1.literal)
        if TwelveClockIdentHash[$2.literal] {
            hour = hour + 12
        }
        $$ = time.Duration(int64(time.Hour) * int64(hour))
    }

%%

type Lexer struct {
    scanner.Scanner
    result Expression
}

var MonthIdentHash = map[string]time.Month{
    "Jan": time.January,
    "Feb": time.February,
    "Mar": time.March,
    "Apr": time.April,
    "May": time.May,
    "June": time.June,
    "July": time.July,
    "Aug": time.August,
    "Sept": time.September,
    "Oct": time.October,
    "Nov": time.November,
    "Dec": time.December,

    "January": time.January,
    "February": time.February,
    "March": time.March,
    "April": time.April,
    "August": time.August,
    "September": time.September,
    "October": time.October,
    "November": time.November,
    "December": time.December,
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
