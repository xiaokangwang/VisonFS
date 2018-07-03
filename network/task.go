package network

type NetworkTaskQueue struct {
}

type NetworkUploadTask struct {
	Filename string
	content  []byte
}
type NetworkUploadTaskResult struct {
}
type NetworkDownloadTask struct {
	Filename string
}
type NetworkDownloadTaskResult struct {
	content []byte
}
type NetworkListTask struct {
	Dir string
}
type NetworkListTaskResult struct {
	Files []string
}

func (ntq *NetworkTaskQueue) EnqueueUploadTask(task NetworkUploadTask)   {}
func (ntq *NetworkTaskQueue) EnqueueDownloadTask(task NetworkUploadTask) {}
func (ntq *NetworkTaskQueue) EnqueueListTask(task NetworkListTask)       {}
