package journeldb

import (
	"io"
	"strconv"
	"strings"
	"sync"
)

type JournelDB struct {
	currentRev uint64
	lock       sync.RWMutex
	canWrite   bool
	jdb        journel
}

func (jdb *JournelDB) StartWrite() {
	jdb.lock.Lock()
	jdb.canWrite = true
	jdb.currentRev++
	jdb.jdb.StartWrite()
}
func (jdb *JournelDB) EndWrite() {
	jdb.canWrite = false
	jdb.lock.Unlock()
}

func (jdb *JournelDB) WriteValue(name, value string) {
	if !jdb.canWrite {
		panic(nil)
	}
	jdb.jdb.WriteValue(name, value)
	if strings.HasPrefix(name, "local") {
		return
	}
	jdb.jdb.ldb.WriteValue(name, value)
}
func (jdb *JournelDB) WriteValueDelete(name string) {
	if !jdb.canWrite {
		panic(nil)
	}
	jdb.jdb.WriteValueDelete(name)
	if strings.HasPrefix(name, "local") {
		return
	}
	jdb.jdb.ldb.WriteValueDelete(name)
}
func (jdb *JournelDB) WriteListValueAdd(name, element string) {
	if !jdb.canWrite {
		panic(nil)
	}
	jdb.jdb.WriteListValueAdd(name, element)
	if strings.HasPrefix(name, "local") {
		return
	}
	jdb.jdb.ldb.WriteListValueAdd(name, element)
}
func (jdb *JournelDB) WriteListValueRemove(name, element string) {
	if !jdb.canWrite {
		panic(nil)
	}
	jdb.jdb.WriteListValueRemove(name, element)
	if strings.HasPrefix(name, "local") {
		return
	}
	jdb.jdb.ldb.WriteListValueRemove(name, element)
}
func (jdb *JournelDB) WriteValueIncrease(name string) {
	if !jdb.canWrite {
		panic(nil)
	}
	jdb.jdb.WriteValueIncrease(name)
	if strings.HasPrefix(name, "local") {
		return
	}
	jdb.jdb.ldb.WriteValueIncrease(name)
}
func (jdb *JournelDB) WriteValueDecrease(name string) {
	if !jdb.canWrite {
		panic(nil)
	}
	jdb.jdb.WriteValueIncrease(name)
	if strings.HasPrefix(name, "local") {
		return
	}
	jdb.jdb.ldb.WriteValueIncrease(name)
}

func (jdb *JournelDB) GetValue(name string) string {
	return jdb.jdb.ldb.GetValue(name)
}
func (jdb *JournelDB) GetValueList(name string) []string {
	return jdb.jdb.ldb.GetValueList(name)
}
func (jdb *JournelDB) WriteRevID() {
	s := strconv.FormatUint(jdb.currentRev, 10)
	jdb.jdb.ldb.ledb.Put([]byte("rev:current"), []byte(s), nil)
}

func Reproduce(ReaderFunc func(name string) io.Reader, checkpoint string, revlessthan uint64, localdb string) *JournelDB {
	db := Open(localdb)
	rev := db.jdb.Reproduce(ReaderFunc, checkpoint, revlessthan)
	db.currentRev = rev
	return db
}
func (jdb *JournelDB) CreateCheckpoint(wr io.Writer, rev uint64) {
	jdb.jdb.CreateCheckpoint(wr, rev)
}
func (jdb *JournelDB) CreateRev(rev uint64, wr io.Writer) {
	jdb.jdb.CreateRev(rev, wr)
}
func Open(localdb string) *JournelDB {
	d := &JournelDB{}
	ldb := d.localdbOpenDatabase(localdb)
	d.jdb.ldb = ldb
	ld, err := ldb.ledb.Get([]byte("rev:current"), nil)
	if err != nil {
		panic(err)
	}
	io, err := strconv.ParseUint(string(ld), 10, 64)
	if err != nil {
		panic(err)
	}
	d.currentRev = io
	return d
}
