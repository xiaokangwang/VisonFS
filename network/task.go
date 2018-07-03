package network

type NetworkTaskQueue struct {
}
type NetworkUploadTask struct {
}
type NetworkUploadTaskResult struct {
}
type NetworkDownloadTask struct {
}
type NetworkDownloadTaskResult struct {
}
type NetworkListTask struct {
}
type NetworkListTaskResult struct {
}

func (ntq *NetworkTaskQueue) EnqueueUploadTask(task NetworkUploadTask)   {}
func (ntq *NetworkTaskQueue) EnqueueDownloadTask(task NetworkUploadTask) {}
func (ntq *NetworkTaskQueue) EnqueueListTask(task NetworkListTask)       {}
