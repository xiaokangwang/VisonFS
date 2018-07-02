package journeldb

import (
	"io"
	"sync"
)

type JournelDB struct {
	currentRev uint64
	lock       sync.RWMutex
	canWrite   bool
}

func (jdb *JournelDB) StartWrite() {}
func (jdb *JournelDB) EndWrite()   {}

func (jdb *JournelDB) WriteValue(name, value string)             {}
func (jdb *JournelDB) WriteValueDelete(name string)              {}
func (jdb *JournelDB) WriteListValueAdd(name, element string)    {}
func (jdb *JournelDB) WriteListValueRemove(name, element string) {}
func (jdb *JournelDB) WriteValueIncrease(name string)            {}
func (jdb *JournelDB) WriteValueDecrease(name string)            {}

func (jdb *JournelDB) GetValue(name string) string {
	return ""
}
func (jdb *JournelDB) GetValueList(name string) []string {
	return nil
}

func Reproduce(ReaderFunc func(name string) io.Reader, checkpoint string, revlessthan uint64, localdb string) *JournelDB {
	return nil
}
func Open(localdb string) *JournelDB {
	return nil
}
