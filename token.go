package logfmt

type token int

const (
	tEOF token = iota - 1
	tError
	tEqual
	tString
	tNumber
	tIdent
)
