package skip

import "testing"

func ExtraSigs() {

}

type Obj struct {
}

func (o Obj) ExtraSigs() {}

func TestSkipExtraSigsFunction(t *testing.T) {
	ExtraSigs()
}

func TestSkipExtraSigsMethod(t *testing.T) {
	Obj{}.ExtraSigs()
}

func TestSkipNoExtraSigs(t *testing.T) { // want "Function TestSkipNoExtraSigs missing the call to method parallel"
}
