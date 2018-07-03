package sync

type PendingSync struct {
	cacheDir          string
	cacheusing        uint64
	cacheCap          uint64
	LastUploadMetaRev uint64
}

func (ps *PendingSync) BlobUpload(content []byte) string {

}
func (ps *PendingSync) BlobGet(hash string) []byte {

}
func (ps *PendingSync) SyncMeta()                       {}
func (ps *PendingSync) CheckoutMeta(revlessthan uint64) {}
func (ps *PendingSync) CreateCheckpoint()               {}
