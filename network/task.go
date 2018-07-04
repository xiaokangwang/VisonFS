package network

type NetworkTaskQueue struct {
}

type NetworkUploadTask struct {
	Filename string
	Content  []byte
}
type NetworkUploadTaskResult struct {
}
type NetworkDownloadTask struct {
	Filename string
}
type NetworkDownloadTaskResult struct {
	Content []byte
}
type NetworkListTask struct {
	Dir string
}
type NetworkListTaskResult struct {
	Files []string
}

func (ntq *NetworkTaskQueue) EnqueueUploadTask(task NetworkUploadTask) {}
func (ntq *NetworkTaskQueue) EnqueueDownloadTask(task NetworkDownloadTask) NetworkDownloadTaskResult {
	panic(nil)
}
func (ntq *NetworkTaskQueue) EnqueueListTask(task NetworkListTask) NetworkListTaskResult {
	panic(nil)
}
