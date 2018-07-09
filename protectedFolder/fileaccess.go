package protectedFolder

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/sha3"

	"github.com/xiaokangwang/VisonFS/transform"
)

type DelegatedAccess struct {
	root string
	key  string
	tf   *transform.Transform
}

const keySize = 32
const nonceSize = 24

func NewDelegatedAccess(tf *transform.Transform, root, key string) *DelegatedAccess {
	return &DelegatedAccess{tf: tf, root: root + "/autocommit", key: key}
}

func (da *DelegatedAccess) ReadToken(t string) string {
	key := da.key
	ssep := strings.Split(t, ".")
	bd := base64.NewDecoder(base64.RawURLEncoding, strings.NewReader(ssep[0]))
	bnd := base64.NewDecoder(base64.RawURLEncoding, strings.NewReader(ssep[1]))
	bdb, _ := ioutil.ReadAll(bd)
	bndb, _ := ioutil.ReadAll(bnd)
	nonce := new([nonceSize]byte)
	copy(nonce[:], bndb[:nonceSize])
	pw := sha3.Sum256([]byte(key))
	res, ok := secretbox.Open(nil, bdb, nonce, &pw)
	if !ok {
		return ""
	}
	return string(res)
}

func (da *DelegatedAccess) CreateToken(c string) string {
	key := da.key
	toenc := []byte(c)
	nonce := new([nonceSize]byte)
	// Read bytes from random and put them in nonce until it is full.
	io.ReadFull(rand.Reader, nonce[:])
	pw := sha3.Sum256([]byte(key))
	ops := secretbox.Seal(nil, toenc, nonce, &pw)
	var mb, dab bytes.Buffer
	mb64e := base64.NewEncoder(base64.RawURLEncoding, &mb)
	mb64e.Write(nonce[:])
	mb64e.Close()
	d64e := base64.NewEncoder(base64.RawURLEncoding, &dab)
	d64e.Write(ops[:])
	d64e.Close()
	return string(dab.Bytes()) + "." + string(mb.Bytes())
}
func (da *DelegatedAccess) ReadFile(path string) ([]byte, error) {
	pn, fn := da.toPath(path)
	dirv := da.root + "/" + pn
	err := os.MkdirAll(dirv, 0700)
	if err != nil {
		return nil, err
	}
	fc, err := ioutil.ReadFile(dirv + "/" + fn)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	da.tf.Decrypt(&buf, bytes.NewReader(fc))
	println(path, string(buf.Bytes()))
	return buf.Bytes(), err
}
func (da *DelegatedAccess) toPath(path string) (string, string) {
	fmt.Println(path)
	dir := strings.Split(path, "/")
	knowndir := "/"
	skip_search := false
dir_for:
	for k := range dir {
		if !skip_search {
			res, err := da.listFileE(knowndir)
			if err != nil {
				panic(err)
			}
			for _, resi := range res {
				fmt.Println("Comparing:", resi.Name(), dir[k])
				if resi.Name() == dir[k] {
					dir[k] = resi.(*decryptFileinfo).inner.Name()
					knowndir += resi.(*decryptFileinfo).inner.Name()
					knowndir += "/"
					continue dir_for
				}
			}
			skip_search = true
		}
		dir[k] = da.CreateToken(dir[k])
	}
	fn := dir[len(dir)-1]
	pn := strings.Join(dir[:len(dir)-1], "/")
	fmt.Println(dir)
	return pn, fn
}
func (da *DelegatedAccess) WriteFile(path string, filecont []byte) error {
	pn, fn := da.toPath(path)
	dirv := da.root + "/" + pn
	os.MkdirAll(dirv, 0700)
	f, err := os.Create(dirv + "/" + fn)
	if err != nil {
		panic(err)
		//return err
	}
	da.tf.Encrypt(f, bytes.NewReader(filecont))
	f.Close()
	return nil
}
func (da *DelegatedAccess) ListFile(path string) ([]os.FileInfo, error) {
	pn, fn := da.toPath(path)
	dirv := da.root + "/" + pn
	dir := dirv + "/" + fn
	if path == "" {
		dir = da.root + "/"
	}
	fmt.Printf("\n\nLIST %v\n", dirv)
	var fni []os.FileInfo
	dird, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, v := range dird {
		fni = append(fni, &decryptFileinfo{inner: v, da: da})
	}

	/*
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				fmt.Printf("\n\nNOINFO %v\n", path)
				return nil
			}
			fmt.Printf("\n\n%v Dir: %v\n\n", info.Name(), info.IsDir())
			fni = append(fni, &decryptFileinfo{inner: info, da: da})
			return nil
		})*/
	return fni, nil
}
func (da *DelegatedAccess) listFileE(epath string) ([]os.FileInfo, error) {
	dirv := da.root + epath
	var fni []os.FileInfo
	dird, err := ioutil.ReadDir(dirv)
	fmt.Println("dirv:"+dirv, err)
	for _, v := range dird {
		fni = append(fni, &decryptFileinfo{inner: v, da: da})
	}

	return fni, nil
}
func (da *DelegatedAccess) RemoveFile(path string) error {
	pn, fn := da.toPath(path)
	dirv := da.root + "/" + pn
	err := os.Remove(dirv + "/" + fn)
	if err != nil {
		return err
	}
	return nil
}

func (da *DelegatedAccess) FileAttr(path string) (os.FileInfo, error) {
	pn, fn := da.toPath(path)
	dirv := da.root + "/" + pn
	info, err := os.Stat(dirv + "/" + fn)
	if err != nil {
		return nil, err
	}
	return &decryptFileinfo{inner: info, da: da}, nil
}

type decryptFileinfo struct {
	inner os.FileInfo
	da    *DelegatedAccess
}

func (df *decryptFileinfo) IsDir() bool        { return df.inner.IsDir() }
func (df *decryptFileinfo) ModTime() time.Time { return df.inner.ModTime() }
func (df *decryptFileinfo) Mode() os.FileMode  { return df.inner.Mode() }
func (df *decryptFileinfo) Name() string {
	fmt.Println(df)
	fmt.Println(df.da)
	fmt.Println(df.inner)
	fmt.Println("decrypt:" + df.da.ReadToken(df.inner.Name()))
	return df.da.ReadToken(df.inner.Name())
}
func (df *decryptFileinfo) Size() int64      { return df.inner.Size() }
func (df *decryptFileinfo) Sys() interface{} { return df.inner.Sys() }
