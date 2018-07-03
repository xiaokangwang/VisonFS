package transform

import (
	"io"

	"github.com/klauspost/reedsolomon"
)

type rsPass struct {
	datashard   int
	parityshard int
	enc         reedsolomon.Encoder
}

func (rp *rsPass) ensureRs() {
	if rp.enc == nil {
		var err error
		rp.enc, err = reedsolomon.New(rp.datashard, rp.parityshard)
		if err != nil {
			panic(err)
		}
	}
}
func (rp *rsPass) PassForword(f []byte) [][]byte {
	rp.ensureRs()
	ot, err := rp.enc.Split(f)
	if err != nil {
		panic(err)
	}
	e := rp.enc.Encode(ot)
	if e != nil {
		panic(err)
	}
	return ot
}
func (rp *rsPass) PassReverse(b [][]byte, len int, w io.Writer) {
	rp.ensureRs()
	err := rp.enc.ReconstructData(b)
	if err != nil {
		panic(err)
	}
	err = rp.enc.Join(w, b, len)
	if err != nil {
		panic(err)
	}
}
