%{
package main

import (
    "text/scanner"
    "strings"
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

type PriorityExpr struct {
    priority int
}

%}

%union{
    token Token
    expr Expression
}

%type<expr> filter
%type<expr> expr
%token<token> STRING
%token<token> NUMBER
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
%%

type Lexer struct {
    scanner.Scanner
    result Expression
}

func (l *Lexer) Lex(lval *yySymType) int {
    token := int(l.Scan())
    switch token {
        case scanner.Ident: token = STRING
        case scanner.Int:   token = NUMBER
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
