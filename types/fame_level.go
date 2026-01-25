package types

type FameLevel string

const (
	FameLevelObscure      FameLevel = "Obscure"
	FameLevelRumored      FameLevel = "Rumored"
	FameLevelNoted        FameLevel = "Noted"
	FameLevelRecognized   FameLevel = "Recognized"
	FameLevelDistinguished FameLevel = "Distinguished"
	FameLevelCelebrated   FameLevel = "Celebrated"
	FameLevelEminent      FameLevel = "Eminent"
	FameLevelRenowned     FameLevel = "Renowned"
	FameLevelIllustrious  FameLevel = "Illustrious"
)

func (f FameLevel) String() string {
	return string(f)
}

func (f FameLevel) Order() int {
	switch f {
	case FameLevelObscure:
		return 1
	case FameLevelRumored:
		return 2
	case FameLevelNoted:
		return 3
	case FameLevelRecognized:
		return 4
	case FameLevelDistinguished:
		return 5
	case FameLevelCelebrated:
		return 6
	case FameLevelEminent:
		return 7
	case FameLevelRenowned:
		return 8
	case FameLevelIllustrious:
		return 9
	}
	return 0
}