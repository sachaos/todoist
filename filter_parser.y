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

const (
    DUE_ON int = iota
    DUE_BEFORE
    DUE_AFTER
)

type DueDateExpr struct {
    operation int
    datetime time.Time
    allDay bool
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
%type<expr> s_overdue
%type<expr> s_time
%token<token> STRING NUMBER
%token<token> MONTH_IDENT TWELVE_CLOCK_IDENT
%token<token> TODAY_IDENT TOMORROW_IDENT YESTERDAY_IDENT
%token<token> DUE BEFORE AFTER OVER OVERDUE
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
    | s_overdue
    {
        $$ = DueDateExpr{allDay: false, datetime: now(), operation: DUE_BEFORE}
    }
    | DUE BEFORE ':' s_datetime
    {
        e := $4.(DueDateExpr)
        e.operation = DUE_BEFORE
        $$ = e
    }
    | DUE AFTER ':' s_datetime
    {
        e := $4.(DueDateExpr)
        e.operation = DUE_AFTER
        $$ = e
    }
    | s_datetime

s_overdue
    : OVER DUE
    {
        $$ = $1
    }
    | OVERDUE
    {
        $$ = $1
    }

s_datetime
    : s_date_year s_time
    {
        date := $1.(time.Time)
        time := $2.(time.Duration)
        $$ = DueDateExpr{allDay: false, datetime: date.Add(time)}
    }
    | s_date_year
    {
        $$ = DueDateExpr{allDay: true, datetime: $1.(time.Time)}
    }
    | s_time
    {
        nd := now().Sub(today())
        d := $1.(time.Duration)
        if (d <= nd) {
          d = d + time.Duration(int64(time.Hour) * 24)
        }
        $$ = DueDateExpr{allDay: false, datetime: today().Add(d)}
    }

s_date_year
    : NUMBER '/' NUMBER '/' NUMBER
    {
        $$ = time.Date(atoi($5.literal), time.Month(atoi($1.literal)), atoi($3.literal), 0, 0, 0, 0, time.Local)
    }
    | MONTH_IDENT NUMBER NUMBER
    {
        $$ = time.Date(atoi($3.literal), MonthIdentHash[strings.ToLower($1.literal)], atoi($2.literal), 0, 0, 0, 0, time.Local)
    }
    | NUMBER MONTH_IDENT NUMBER
    {
        $$ = time.Date(atoi($3.literal), MonthIdentHash[strings.ToLower($2.literal)], atoi($1.literal), 0, 0, 0, 0, time.Local)
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
    | TODAY_IDENT
    {
        $$ = today()
    }
    | TOMORROW_IDENT
    {
        $$ = today().AddDate(0, 0, 1)
    }
    | YESTERDAY_IDENT
    {
        $$ = today().AddDate(0, 0, -1)
    }

s_date
    : MONTH_IDENT NUMBER
    {
        $$ = time.Date(today().Year(), MonthIdentHash[strings.ToLower($1.literal)], atoi($2.literal), 0, 0, 0, 0, time.Local)
    }
    | NUMBER MONTH_IDENT
    {
        $$ = time.Date(today().Year(), MonthIdentHash[strings.ToLower($2.literal)], atoi($1.literal), 0, 0, 0, 0, time.Local)
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
    "jan": time.January,
    "feb": time.February,
    "mar": time.March,
    "apr": time.April,
    "may": time.May,
    "june": time.June,
    "july": time.July,
    "aug": time.August,
    "sept": time.September,
    "oct": time.October,
    "nov": time.November,
    "dec": time.December,

    "january": time.January,
    "february": time.February,
    "march": time.March,
    "april": time.April,
    "august": time.August,
    "september": time.September,
    "october": time.October,
    "november": time.November,
    "december": time.December,
}

var TwelveClockIdentHash = map[string]bool{
    "am": false,
    "pm": true,
}

var TodayIdentHash = map[string]bool {
    "today": true,
    "tod": true,
}

var OverDueHash = map[string]bool {
    "overdue": true,
    "od": true,
}


func (l *Lexer) Lex(lval *yySymType) int {
    token := int(l.Scan())
    switch token {
        case scanner.Ident:
            lowerToken := strings.ToLower(l.TokenText())
            if _, ok := MonthIdentHash[lowerToken]; ok {
                token = MONTH_IDENT
            } else if _, ok := TwelveClockIdentHash[lowerToken]; ok {
                token = TWELVE_CLOCK_IDENT
            } else if _, ok := TodayIdentHash[lowerToken]; ok {
                token = TODAY_IDENT
            } else if lowerToken == "tomorrow" {
                token = TOMORROW_IDENT
            } else if lowerToken == "yesterday" {
                token = YESTERDAY_IDENT
            } else if lowerToken == "due" {
                token = DUE
            } else if lowerToken == "before" {
                token = BEFORE
            } else if lowerToken == "after" {
                token = AFTER
            } else if lowerToken == "over" {
                token = OVER
            } else if _, ok := OverDueHash[lowerToken]; ok {
                token = OVERDUE
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
