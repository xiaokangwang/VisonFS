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

func (ins *Instance) Prepare(gitpath, pubdir, prvdir, prvpass, tmpdir, uploadperfix string) {
	//instance all dep
	ins.gitctli = gitctl.NewGitctl(gitpath)
	ins.transformi = transform.NewTransform(pubdir, prvdir, prvpass)
	ins.networki = network.NewNetworkTaskQueue(uploadperfix)
	ins.synci = sync.NewPendingSync(ins.transformi, ins.networki, tmpdir)
	ins.protectedFolderi = protectedFolder.NewDelegatedAccess(ins.transformi, gitpath, prvpass)
	ins.filei = file.NewFileTree(ins.transformi, ins.protectedFolderi, ins.synci, ins.gitctli)
	//Look for dirty
	ins.synci.UploadDirty()
}
func (ins *Instance) Launch() *file.FileTree {
	return ins.filei
}
func (ins *Instance) Purge() {
	ins.synci.Purge()
}
