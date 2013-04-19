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
		switch sc.ss {
		case stepBeginKey:
			mark := sc.off - 1
			sc.scanWhile(stepContinue)
			return scanKey, sc.b[mark : sc.off-1]
		case stepBeginValue:
			mark := sc.off - 1
			sc.scanWhile(stepContinue)
			return scanVal, sc.b[mark : sc.off-1]
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
	for sc.off < len(sc.b) {
		sc.ss = sc.s.step(sc.b[sc.off])
		sc.off++
		if sc.ss != what {
			return
		}
	}
	if sc.off == len(sc.b) {
		sc.off++
	}
	sc.ss = stepEnd
}
