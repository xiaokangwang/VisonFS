package journeldb

import "io"

type journel struct {
	ldb *localDb
}

func (jdb *journel) Reproduce(ReaderFunc func(name string) io.Reader, checkpoint string, revlessthan uint64) uint64 {
	return 0
}

func (jdb *journel) CreateCheckpoint(wr io.Writer)      {}
func (jdb *journel) CreateRev(rev uint64, wr io.Writer) {}

func (jdb *journel) WriteValue(name, value string)             {}
func (jdb *journel) WriteValueDelete(name string)              {}
func (jdb *journel) WriteListValueAdd(name, element string)    {}
func (jdb *journel) WriteListValueRemove(name, element string) {}
func (jdb *journel) WriteValueIncrease(name string)            {}
func (jdb *journel) WriteValueDecrease(name string)            {}

func (jdb *journel) StartWrite()         {}
func (jdb *journel) EndWrite(Rev uint64) {}
