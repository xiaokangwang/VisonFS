package sync

import "github.com/xiaokangwang/VisonFS/journeldb"
import "github.com/xiaokangwang/VisonFS/transform"

type PendingSync struct {
	cacheDir          string
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
func (ps *PendingSync) SyncMeta()                       {}
func (ps *PendingSync) CheckoutMeta(revlessthan uint64) {}
func (ps *PendingSync) CreateCheckpoint()               {}
