package logfmt

import (
	"io"
)

// A Decoder reads and decodes JSON objects from an input stream.
type Decoder struct {
	r    io.Reader
	buf  []byte
	scan scanner
	err  error
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (dec *Decoder) Decode(v interface{}) error {
	if dec.err != nil {
		return dec.err
	}

	n, err := dec.readValue()
	if err != nil {
		return err
	}

	err = unmarshal(dec.buf[0:n], v)

	rest := copy(dec.buf, dec.buf[n:])
	dec.buf = dec.buf[0:rest]

	return err
}

// readValue reads a JSON value into dec.buf.
// It returns the length of the encoding.
func (dec *Decoder) readValue() (int, error) {
	var err error
	var scanp int
	dec.scan.reset()
	for {
		for i, c := range dec.buf[scanp:] {
			v := dec.scan.step(&dec.scan, rune(c))
			switch v {
			case scanEnd:
				scanp += i + 1
				return scanp, nil
			case scanError:
				dec.err = dec.scan.err
				return 0, dec.err
			}
		}
		scanp = len(dec.buf)
		if err != nil {
			if err == io.EOF {
				if nonSpace(dec.buf) {
					err = io.ErrUnexpectedEOF
				}
			}
			dec.err = err
			return 0, dec.err
		}
		const minRead = 512
		if cap(dec.buf)-len(dec.buf) < minRead {
			newBuf := make([]byte, len(dec.buf), 2*cap(dec.buf)+minRead)
			copy(newBuf, dec.buf)
			dec.buf = newBuf
		}
		var n int
		n, err = dec.r.Read(dec.buf[len(dec.buf):cap(dec.buf)])
		dec.buf = dec.buf[0 : len(dec.buf)+n]
	}
	return scanp, nil
}

func nonSpace(buf []byte) bool {
	for _, c := range buf {
		switch c {
		case ' ', '\t', '\r', '\n':
		default:
			return true
		}
	}
	return false
}
