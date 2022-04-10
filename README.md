# suela

## Introduction

SUELA is a Simple URL Embeddable Language - Implemented in Go.

It is designed to be able to use it in the query part of http or https
URLs without too much escaping. As such SUELA does not allow whitespace,
or comments. Unicode is supported in string values but that will have to be
URLescaped. The language itself only uses these url safe characters:
A-Za-z0-9-._~()'!*:@,;/

SUELA is also designed so it can be used for filtering JSON documents
that are fetched from some underlying data store.



