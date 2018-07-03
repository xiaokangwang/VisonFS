package sync

import (
	"bytes"
	"io"
	"os"
	"strconv"

	"github.com/xiaokangwang/VisonFS/cache"
	"github.com/xiaokangwang/VisonFS/journeldb"
	"github.com/xiaokangwang/VisonFS/transform"
)

type PendingSync struct {
	cacheDir          string
	metadomain        string
	cacheusing        uint64
	cacheCap          uint64
	LastUploadMetaRev uint64
	jd                *journeldb.JournelDB
	tf                *transform.Transform
}

func (ps *PendingSync) BlobUpload(content []byte) string {

}
func (ps *PendingSync) BlobGet(hash string) []byte {

}
func (ps *PendingSync) SyncMeta() {
	syncto := ps.jd.CurrentRev()

	uploading := ps.LastUploadMetaRev + 1
	uploadingFirst := uploading
	var syncbuf bytes.Buffer
	for uploading <= syncto {
		ps.jd.CreateRev(uploading, &syncbuf)
		uploading++
	}

	fname := ps.cacheDir + "/meta/" + ps.metadomain + "_rev_" + strconv.FormatUint(uploadingFirst, 10)
	ps.QueueFileNetworkUpload(fname, syncbuf.Bytes())

}
func (ps *PendingSync) QueueFileNetworkUpload(fname string, content []byte) {
	cache.SetDirty(fname)
	f, err := os.Create(fname)
	io.Copy(f, bytes.NewBuffer(content))
	//TODO:Queue Upload
	cache.RemoveDirty(fname)
}
func (ps *PendingSync) CheckoutMeta(revlessthan uint64) {}
func (ps *PendingSync) CreateCheckpoint()               {}
