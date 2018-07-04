package sync

import (
	"bytes"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/xiaokangwang/VisonFS/cache"
	"github.com/xiaokangwang/VisonFS/network"
	"github.com/xiaokangwang/VisonFS/transform"
	"golang.org/x/crypto/sha3"
)

type PendingSync struct {
	cacheDir   string
	metadomain string
	cacheusing uint64
	cacheCap   uint64

	//jd                *journeldb.JournelDB
	tf   *transform.Transform
	nw   *network.NetworkTaskQueue
	crlo sync.RWMutex
}

func NewPendingSync(tf *transform.Transform, nw *network.NetworkTaskQueue) *PendingSync {
	return &PendingSync{tf: tf, nw: nw}
}

func (ps *PendingSync) BlobUpload(content []byte) string {
	//transform
	out, cookie := ps.tf.Advance(content)
	syncookie := cookie
	for n := range out {
		sum := sha3.Sum256(out[n])
		sumx := hex.EncodeToString(sum[:])
		ps.QueueFileNetworkUpload("blob/"+sumx, content)
		syncookie += "$"
		syncookie += sumx
	}
	return syncookie

}
func (ps *PendingSync) BlobGet(hash string) []byte {
	cookie := strings.Split(hash, "$")
	cookiei := cookie[0]
	var file [][]byte
	for k := range cookie {
		if k == 0 {
			cookiei = cookie[0]
		} else {
			fc, err := ps.QueueFileNetworkDownload(cookie[k])
			if err != nil {
				panic(err)
			}
			file = append(file, fc)
		}
	}
	return ps.tf.Reverse(file, cookiei)
}

/*
func (ps *PendingSync) SyncMeta() {
	syncto := ps.jd.CurrentRev()

	uploading := ps.LastUploadMetaRev + 1
	uploadingFirst := uploading
	var syncbuf bytes.Buffer
	for uploading <= syncto {
		ps.jd.CreateRev(uploading, &syncbuf)
		uploading++
	}
	var syncbufenc bytes.Buffer
	by := syncbuf.Bytes()
	syncbufr := bytes.NewReader(by)
	ps.tf.Encrypt(&syncbufenc, syncbufr)

	fname := ps.cacheDir + "/meta/" + ps.metadomain + "_rev_" + strconv.FormatUint(uploadingFirst, 10)
	ps.QueueFileNetworkUpload(fname, syncbufenc.Bytes())

}*/
func (ps *PendingSync) QueueFileNetworkUpload(fname string, content []byte) {
	cache.SetDirty(ps.cacheDir + "/" + fname)
	ps.crlo.Lock()
	f, err := os.Create(ps.cacheDir + "/" + fname)
	if err != nil {
		panic(err)
	}
	io.Copy(f, bytes.NewBuffer(content))
	f.Close()
	ps.crlo.Unlock()
	//TODO:Queue Upload
	var dt network.NetworkUploadTask
	dt.Filename = fname
	dt.Content = content
	ps.nw.EnqueueUploadTask(dt)
	cache.RemoveDirty(ps.cacheDir + "/" + fname)
}
func (ps *PendingSync) QueueFileNetworkDownload(fname string) ([]byte, error) {
	if cache.IsExist(ps.cacheDir + "/" + fname) {
		ps.crlo.RLock()
		c, e := ioutil.ReadFile(ps.cacheDir + "/" + fname)
		ps.crlo.RUnlock()
		return c, e
	}
	//TODO：DownloadFile
	var dt network.NetworkDownloadTask
	dt.Filename = fname
	ou := ps.nw.EnqueueDownloadTask(dt)
	//Write cache
	ps.crlo.Lock()
	ps.crlo.Unlock()
	return ou.Content, nil
}
func (ps *PendingSync) UploadDirty() {
	res := cache.FindDrity(ps.cacheDir)
	for _, v := range res {
		c, _ := ioutil.ReadFile(ps.cacheDir + "/" + v)
		var dt network.NetworkUploadTask
		dt.Filename = v
		dt.Content = c
		ps.nw.EnqueueUploadTask(dt)
		cache.RemoveDirty(ps.cacheDir + "/" + v)
	}
}
func (ps *PendingSync) Purge() {
	cache.Purge(ps.cacheDir)
}
