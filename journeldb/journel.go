package journeldb

import (
	"bufio"
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
	var reader io.Reader
	if checkpoint != "" {
		reader = ReaderFunc(checkpoint)
	}
	var rev uint64
	for {
		if reader == nil {
			s := strconv.FormatUint(rev, 10)
			reader = ReaderFunc(s)
		}
		if reader == nil {
			break
		}
		linereader := bufio.NewReader(reader)
		for {
			reading, err := linereader.ReadString('\n')
			if err != nil {
				break
			}
			if strings.HasPrefix(reading, "#REV") {
				stlo := strings.Split(reading, " ")
				io, err := strconv.ParseUint(stlo[1], 10, 64)
				if err != nil {
					panic(err)
				}
				rev = io
				if rev > revlessthan && revlessthan != 0 {
					return rev
				}
			}
			jdb.reproduceline(reading)
		}

	}
	return rev
}
func (jdb *journel) reproduceline(obj string) {
	stlo := strings.Split(obj, " ")
	stlo[1], _ = url.QueryUnescape(stlo[1])
	if len(stlo) == 3 {
		stlo[2], _ = url.QueryUnescape(stlo[2])
	}

	switch stlo[0] {
	case "WriteValue":
		jdb.ldb.WriteValue(stlo[1], stlo[2])
	case "WriteValueDelete":
		jdb.ldb.WriteValueDelete(stlo[1])
	case "WriteListValueAdd":
		jdb.ldb.WriteListValueAdd(stlo[1], stlo[2])
	case "WriteListValueRemove":
		jdb.ldb.WriteListValueRemove(stlo[1], stlo[2])
	case "WriteValueIncrease":
		jdb.ldb.WriteValueIncrease(stlo[1])
	case "WriteValueDecrease":
		jdb.ldb.WriteValueDecrease(stlo[1])
	}
}

func (jdb *journel) CreateCheckpoint(wr io.Writer, rev uint64) {
	i := jdb.ldb.ledb.NewIterator(nil, nil)
	for i.Next() {
		if strings.HasPrefix(string(i.Key()), "rev:") {
			continue
		}
		fmt.Fprintf(wr, "WriteValue %s %s\n", string(i.Key()), string(i.Value()))
	}
	s := strconv.FormatUint(rev, 10)
	fmt.Fprintf(wr, "#REVCHECKPOINT %s\n", s)
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
