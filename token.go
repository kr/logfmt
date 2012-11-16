package logfmt

const (
	tEOF = iota - 1
	tError
	tEqual
	tString
	tNumber
	tIdent
)

type token struct {
	t   int
	src []byte
}

func (tok *token) isEOF() bool {
	return tok.t == tEOF
}

func (tok *token) isKey() bool {
	return tok.t == tString || tok.t == tIdent
}

func (tok *token) isVal() bool {
	return tok.t == tString || tok.t == tIdent || tok.t == tNumber
}

func (tok *token) isEqual() bool {
	return tok.t == tEqual
}
