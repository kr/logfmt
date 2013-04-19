package logfmt

type scannerType int

const (
	scanKey scannerType = iota
	scanVal
	scanEnd
)

type scanner struct {
	s   *stepper
	b   []byte
	off int
	ss  stepperState
}

func newScanner(b []byte) *scanner {
	return &scanner{s: new(stepper), ss: stepSkip}
}

func (sc *scanner) next() (scannerType, []byte) {
	for {
		switch sc.ss {
		case stepBeginKey:
			mark := sc.off
			sc.scanWhile(stepContinue)
			return scanKey, sc.b[mark:sc.off]
		case stepBeginValue:
			mark := sc.off
			sc.scanWhile(stepContinue)
			return scanVal, sc.b[mark:sc.off]
		case stepEnd:
			return scanEnd, nil
		default:
			sc.scanWhile(stepSkip)
		}
	}
}

func (sc *scanner) scanWhile(what stepperState) {
	for _, c := range sc.b[sc.off:] {
		sc.off++
		if sc.ss = sc.s.step(c); sc.ss != what {
			return
		}
	}
}
