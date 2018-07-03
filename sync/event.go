package sync

type PendingSync struct{}

func (ps *PendingSync) BlobUpload()       {}
func (ps *PendingSync) BlobGet()          {}
func (ps *PendingSync) SyncMeta()         {}
func (ps *PendingSync) CreateCheckpoint() {}
