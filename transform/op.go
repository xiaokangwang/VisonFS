package transform

type Transform struct {
	GtTCookie string
	LtTCookie string
}

/*Transfrom Steps:
gt 4MB
Spilt -> RS encoding -> gpg encryption
lt 4MB
Dup -> gpg encryption
*/

const Threshold = 1024 * 1024 * 4

func (t *Transform) Advance(f []byte) ([][]byte, string) {
	transformMethod := "???"
	if len(f) > Threshold {
		transformMethod = t.GtTCookie
	} else {
		transformMethod = t.LtTCookie
	}
	panic(nil)
}
func (t *Transform) Reverse(b [][]byte, c string) []byte {
	panic(nil)
}
func (t *Transform) NeedAtLeast(c string) int {
	panic(nil)
}
