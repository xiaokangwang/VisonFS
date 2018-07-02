package journeldb

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type journel struct {
	ldb *localDb
	buf *bytes.Buffer
}

func (jdb *journel) Reproduce(ReaderFunc func(name string) io.Reader, checkpoint string, revlessthan uint64) uint64 {
	return 0
}

func (jdb *journel) CreateCheckpoint(wr io.Writer) {
	i := jdb.ldb.ledb.NewIterator(nil, nil)
	for i.Next() {
		if strings.HasPrefix(string(i.Key()), "rev:") {
			continue
		}
		fmt.Fprintf(wr, "WriteValue %s %s\n", string(i.Key()), string(i.Value()))
	}
}
func (jdb *journel) CreateRev(rev uint64, wr io.Writer) {
	s := strconv.FormatUint(rev, 10)
	so := jdb.ldb.GetValue("rev:" + s)
	fmt.Fprintf(wr, "#REVSTART %s\n", s)
	wr.Write([]byte(so))
	fmt.Fprintf(wr, "#REVEND %s\n", s)
}

func (jdb *journel) WriteValue(name, value string) {
	name = url.QueryEscape(name)
	value = url.QueryEscape(value)
	fmt.Fprintf(jdb.buf, "WriteValue %s %s\n", name, value)
}
func (jdb *journel) WriteValueDelete(name string) {
	name = url.QueryEscape(name)
	fmt.Fprintf(jdb.buf, "WriteValueDelete %s\n", name)
}
func (jdb *journel) WriteListValueAdd(name, element string) {
	name = url.QueryEscape(name)
	element = url.QueryEscape(element)
	fmt.Fprintf(jdb.buf, "WriteListValueAdd %s %s\n", name, element)
}
func (jdb *journel) WriteListValueRemove(name, element string) {
	name = url.QueryEscape(name)
	element = url.QueryEscape(element)
	fmt.Fprintf(jdb.buf, "WriteListValueRemove %s %s\n", name, element)
}
func (jdb *journel) WriteValueIncrease(name string) {
	name = url.QueryEscape(name)
	fmt.Fprintf(jdb.buf, "WriteValueIncrease %s\n", name)
}
func (jdb *journel) WriteValueDecrease(name string) {
	name = url.QueryEscape(name)
	fmt.Fprintf(jdb.buf, "WriteValueDecrease %s\n", name)
}

func (jdb *journel) StartWrite() {
	jdb.buf = new(bytes.Buffer)
}
func (jdb *journel) EndWrite(Rev uint64) {
	s := strconv.FormatUint(Rev, 10)
	jdb.ldb.WriteValue("rev:"+s, string(jdb.buf.Bytes()))
}
