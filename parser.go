package suela

/*
FILTER -> STATEMENTS .
ATEMENTS -> STATEMENT nl STATEMENTS | .
STATEMENT -> where CALL | join CALL | CALL .
CALL -> funcname op ARGS cp .
ARGS -> ARG comma ARGS | .
ARG -> string | int | float | field | EXPR .
*/

type Parser struct {
	Lexer
}
