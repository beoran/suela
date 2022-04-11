package suela

/*
SCRIPT -> STATEMENTS .
STATEMENTS -> STATEMENT nl STATEMENTS | .
STATEMENT -> comment | where CALL | join CALL | CALL .
CALL -> funcname op ARGS cp .
ARGS -> ARG comma ARGS | .
ARG -> string | int | float | field | EXPR .
*/

type Parser struct {
	Lexer
}
