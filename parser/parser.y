%{
package parser

import (
  "io"
)
%}

%start program

%union {
  val *TokenValue
  res interface{}
}

%token FROM WHERE SELECT

%token <val> ID NUM STR PATH

%token ':' '(' ')' ','

%left '+' '-' '*' DIVOP '%'

%type <res> program transforms transform fromClause whereClause whereClauseOpt 
  selectClause columns column columnName expr pathExpr literalExpr callExpr
  callParamsOpt callParams binaryExpr

%%

program:
  transforms
  {
    $$ = $1
  }

transforms:
  transforms transform
  {
    $$ = append($1.([]Transform), $2.(Transform))
  }
| transform
  {
    $$ = []Transform{$1.(Transform)}
  }

transform:
  fromClause whereClauseOpt selectClause
  {
    $$ = &Transform{
      From:  $1.(PathExpr),
      Where: $2.(Expr),
      Select: $3.([]Column),
    }
  }

fromClause:
  FROM PATH
  {
    $$ = &PathExpr{Path: $2.Path}
  }

whereClause:
  WHERE expr
  {
    $$ = $2
  }

whereClauseOpt:
  whereClause
  {
    $$ = $1
  }
| {
    $$ = nil
  }

selectClause:
  SELECT columns
  {
    $$ = $2
  }

columns:
  columns ',' column
  {
    $$ = append($1.([]Column), $3.(Column))
  }
| column
  {
    $$ = []Column{$1.(Column)}
  }

column:
  columnName expr
  {
    $$ = &Column {
      Name: $1.(string),
      Expr: $2.(Expr),
    }
  }

columnName:
  STR ':'
  {
    $$ = $1.Str
  }
| ID ':'
  {
    $$ = yylex.(*Lexer).Names.NameFromID($1.ID)
  }

expr:
  binaryExpr
| callExpr
| pathExpr
| literalExpr
| '(' expr ')'
  {
    $$ = $2
  }

pathExpr:
  PATH
  {
    $$ = &PathExpr{Path:$1.Path}
  }

literalExpr:
  NUM
  {
    $$ = &LiteralExpr{Val:$1.Num}
  }
| STR
  {
    $$ = &LiteralExpr{Val:$1.Str}
  }

callExpr:
  ID '(' callParamsOpt ')'
  {
    $$ = &CallExpr{
      FuncName: $1.ID,
      Args:     $3.([]Expr),
    }
  }

callParamsOpt:
  callParams
  {
    $$ = $1
  }
| {
    $$ = nil
  }

callParams:
  callParams ',' expr
  {
    $$ = append($1.([]Expr), $3.(Expr))
  }
| expr
  {
    $$ = []Expr{$1.(Expr)}
  }

binaryExpr:
  expr '+' expr
  {
    $$ = &BinaryExpr{
      LHS: $1.(Expr),
      RHS: $3.(Expr),
      Op:  OpAdd,
    }
  }
| expr '-' expr
  {
    $$ = &BinaryExpr{
      LHS: $1.(Expr),
      RHS: $3.(Expr),
      Op:  OpSub,
    }
  }
| expr '*' expr
  {
    $$ = &BinaryExpr{
      LHS: $1.(Expr),
      RHS: $3.(Expr),
      Op:  OpMul,
    }
  }
| expr DIVOP expr
  {
    $$ = &BinaryExpr{
      LHS: $1.(Expr),
      RHS: $3.(Expr),
      Op:  OpDiv,
    }
  }
| expr '%' expr
  {
    $$ = &BinaryExpr{
      LHS: $1.(Expr),
      RHS: $3.(Expr),
      Op:  OpMod,
    }
  }

%%

func Parse(fileName string, reader io.Reader, names *Names) {
  yyErrorVerbose = true
  l := NewLexer(fileName, reader, names)
  yyParse(l)
}
