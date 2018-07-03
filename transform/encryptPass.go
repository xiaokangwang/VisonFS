package transform

import (
	"io"
	"os"

	"golang.org/x/crypto/openpgp"
	//"golang.org/x/crypto/openpgp/packet"
)

type encryptPass struct {
	publickeyPath  string
	privatekeyPath string
	privatekeyPass string
	publickey      openpgp.EntityList
	privateKey     openpgp.EntityList
}

func (ep *encryptPass) EnsurePublickeyEntity() {
	if ep.publickey == nil {
		pubr, err := os.Open(ep.publickeyPath)
		if err != nil {
			panic(err)
		}
		defer pubr.Close()
		//pr := packet.NewReader(pubr)
		ep.publickey, err = openpgp.ReadKeyRing(pubr)
	}

}
func (ep *encryptPass) EnsurePrivatekeyEntity() {
	if ep.privateKey == nil {
		pubr, err := os.Open(ep.privatekeyPath)
		if err != nil {
			panic(err)
		}
		defer pubr.Close()
		//pr := packet.NewReader(pubr)
		ep.privateKey, err = openpgp.ReadKeyRing(pubr)
	}
}

func (ep *encryptPass) PassForword(w io.Writer, r io.Reader) {
	ep.EnsurePublickeyEntity()
	ep.EnsurePrivatekeyEntity()
	pv := ep.privateKey[0]
	pv.PrivateKey.Decrypt([]byte(ep.privatekeyPass))
	for _, v := range pv.Subkeys {
		v.PrivateKey.Decrypt([]byte(ep.privatekeyPass))
	}
	outd, err := openpgp.Encrypt(w, ep.publickey, pv, nil, nil)
	if err != nil {
		panic(err)
	}
	io.Copy(outd, r)
	outd.Close()
}
func (ep *encryptPass) PassReverse(w io.Writer, r io.Reader) {
	ep.EnsurePublickeyEntity()
	ep.EnsurePrivatekeyEntity()
	pv := ep.privateKey[0]
	pv.PrivateKey.Decrypt([]byte(ep.privatekeyPass))
	for _, v := range pv.Subkeys {
		v.PrivateKey.Decrypt([]byte(ep.privatekeyPass))
	}
	de, err := openpgp.ReadMessage(r, ep.privateKey, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return []byte(ep.privatekeyPass), nil
	}, nil)
	if err != nil {
		panic(err)
	}
	io.Copy(w, de.UnverifiedBody)
	if de.SignatureError != nil {
		panic(de.SignatureError)
	}
}
