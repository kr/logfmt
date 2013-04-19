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
	return &scanner{b: b, s: newStepper(), ss: stepSkip}
}

func (sc *scanner) next() (scannerType, []byte) {
	for {
		println("start ss", sc.ss.String())

		switch sc.ss {
		case stepBeginKey:
			mark := sc.off
			sc.scanWhile(stepContinue)
			return scanKey, sc.b[mark:sc.off]
		case stepBeginValue:
			mark := sc.off
			sc.scanWhile(stepContinue)
			return scanVal, sc.b[mark:sc.off]
		case stepEqual:
			sc.scanWhile(stepEqual)
		case stepEnd:
			return scanEnd, nil
		default:
			sc.scanWhile(stepSkip)
		}
	}
}

func (sc *scanner) scanWhile(what stepperState) {
	for ; sc.off < len(sc.b); sc.off++ {
		sc.ss = sc.s.step(sc.b[sc.off])
		println("ss", sc.ss.String())
		if sc.ss != what {
			return
		}
	}
	sc.ss = stepEnd
}
