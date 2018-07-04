package instanceadm

import "github.com/xiaokangwang/VisonFS/gitctl"
import "github.com/xiaokangwang/VisonFS/file"
import "github.com/xiaokangwang/VisonFS/protectedFolder"
import "github.com/xiaokangwang/VisonFS/network"
import "github.com/xiaokangwang/VisonFS/sync"
import "github.com/xiaokangwang/VisonFS/transform"

type Instance struct {
	gitctli          *gitctl.Gitctl
	filei            *file.FileTree
	protectedFolderi *protectedFolder.DelegatedAccess
	networki         *network.NetworkTaskQueue
	synci            *sync.PendingSync
	transformi       *transform.Transform
}

func (ins *Instance) Prepare() {
	//instance all dep
}
func (ins *Instance) Launch() {

}
func (ins *Instance) LoadConfig() {

}
