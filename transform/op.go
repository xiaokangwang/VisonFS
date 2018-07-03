package transform

import (
	"bytes"
	"strconv"
	"strings"
)

type Transform struct {
	GtTCookie string
	LtTCookie string
	rs        rsPass
	ep        encryptPass
}

/*Transfrom Steps:
gt 4MB
Spilt -> RS encoding -> gpg encryption
lt 4MB
Dup -> gpg encryption
*/

const Threshold = 1024 * 1024 * 4

func (t *Transform) Advance(f []byte) ([][]byte, string) {
	var transformMethod string
	if len(f) > Threshold {
		transformMethod = t.GtTCookie
		dlen := len(f)
		transformMethod += ";" + strconv.Itoa(dlen)
		out := t.rs.PassForword(f)
		for k := range out {
			var transBuffer bytes.Buffer
			t.ep.PassForword(&transBuffer, bytes.NewBuffer(out[k]))
			out[k] = transBuffer.Bytes()
		}
		return out, transformMethod
	} else {
		transformMethod = t.LtTCookie
		out := make([][]byte, 4)
		var i = 0
		for i <= 3 {
			var transBuffer bytes.Buffer
			t.ep.PassForword(&transBuffer, bytes.NewBuffer(f))
			out[i] = transBuffer.Bytes()
			i++
		}
		return out, transformMethod
	}
}
func (t *Transform) Reverse(b [][]byte, c string) []byte {
	if c == t.LtTCookie {
		for k := range b {
			if b[k] != nil {
				var transBuffer bytes.Buffer
				t.ep.PassReverse(&transBuffer, bytes.NewBuffer(b[k]))
				return transBuffer.Bytes()
			}
			panic(nil)
		}
	}
	cookieo := strings.Split(c, ";")
	len, _ := strconv.Atoi(cookieo[1])
	for k := range b {
		if b[k] != nil {
			var transBuffer bytes.Buffer
			t.ep.PassReverse(&transBuffer, bytes.NewBuffer(b[k]))
			b[k] = transBuffer.Bytes()
		}
	}
	var outBuffer bytes.Buffer
	t.rs.PassReverse(b, len, &outBuffer)
	return outBuffer.Bytes()
}
func (t *Transform) NeedAtLeast(c string) int {
	return t.rs.datashard
}
