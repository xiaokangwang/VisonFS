package protectedFolder

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

//const key = "rgsf8o1lqjbttgzn08ssqcbooh3y2yfleige"

func (da *DelegatedAccess) ReadToken(t string) string {
	key := da.key
	ssep := strings.Split(t, "_")
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
	return string(dab.Bytes()) + "_" + string(mb.Bytes())
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
	return buf.Bytes(), err
}
func (da *DelegatedAccess) toPath(path string) (string, string) {
	dir := strings.Split(path, "/")
	for k := range dir {
		dir[k] = da.CreateToken(dir[k])
	}
	fn := dir[len(dir)-1]
	pn := strings.Join(dir[:len(dir)-2], "/")
	return pn, fn
}
func (da *DelegatedAccess) WriteFile(path string, filecont []byte) error {
	pn, fn := da.toPath(path)
	dirv := da.root + "/" + pn
	f, err := os.Create(dirv + "/" + fn)
	if err != nil {
		return err
	}
	da.tf.Encrypt(f, bytes.NewReader(filecont))
	f.Close()
	return nil
}
func (da *DelegatedAccess) ListFile(path string) ([]os.FileInfo, error) {
	pn, fn := da.toPath(path)
	dirv := da.root + "/" + pn
	dir := dirv + "/" + fn
	var fni []os.FileInfo
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		fni = append(fni, &decryptFileinfo{inner: info})
		return nil
	})
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

type decryptFileinfo struct {
	inner os.FileInfo
	da    *DelegatedAccess
}

func (df *decryptFileinfo) IsDir() bool        { return df.inner.IsDir() }
func (df *decryptFileinfo) ModTime() time.Time { return df.inner.ModTime() }
func (df *decryptFileinfo) Mode() os.FileMode  { return df.inner.Mode() }
func (df *decryptFileinfo) Name() string {
	return df.da.ReadToken(df.inner.Name())
}
func (df *decryptFileinfo) Size() int64      { return df.inner.Size() }
func (df *decryptFileinfo) Sys() interface{} { return df.inner.Sys() }
