package journeldb

import (
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

func (jdb *JournelDB) localdbOpenDatabase(loc string) *localDb {
	ld, err := leveldb.OpenFile(loc, nil)
	if err != nil {
		panic(err)
	}
	return &localDb{ledb: ld}
}

type localDb struct {
	ledb *leveldb.DB
}

func (jdb *localDb) WriteValue(name, value string) {
	jdb.ledb.Put([]byte(name), []byte(value), nil)
}
func (jdb *localDb) WriteValueDelete(name string) {
	jdb.ledb.Delete([]byte(name), nil)
}
func (jdb *localDb) WriteListValueAdd(name, element string) {
	vlist := jdb.GetValueList(name)
	boolfound := false
	for _, v := range vlist {
		if element == v {
			boolfound = true
		}
	}
	if !boolfound {
		vlist = append(vlist, element)
	}
	out := strings.Join(vlist, "\n")
	jdb.WriteValue(name, out)
}
func (jdb *localDb) WriteListValueRemove(name, element string) {
	vlist := jdb.GetValueList(name)
	var nextvlist []string
	for _, v := range vlist {
		if element == v {
			continue
		}
		nextvlist = append(nextvlist, v)
	}

	out := strings.Join(nextvlist, "\n")
	jdb.WriteValue(name, out)
}
func (jdb *localDb) WriteValueIncrease(name string) {
	v := jdb.GetValue(name)
	if v == "" {
		v = "0"
	}
	vi, _ := strconv.Atoi(v)
	vi++
	v = strconv.Itoa(vi)
	jdb.WriteValue(name, v)
}
func (jdb *localDb) WriteValueDecrease(name string) {
	v := jdb.GetValue(name)
	if v == "" {
		v = "0"
	}
	vi, _ := strconv.Atoi(v)
	v = strconv.Itoa(vi)
	jdb.WriteValue(name, v)
}

func (jdb *localDb) GetValue(name string) string {
	vl, err := jdb.ledb.Get([]byte(name), nil)
	if err != nil {
		return ""
	}
	return string(vl)
}
func (jdb *localDb) GetValueList(name string) []string {
	v := jdb.GetValue(name)
	return strings.Split(v, "\n")
}
