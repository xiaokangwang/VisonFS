package filesystem

import (
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/xiaokangwang/VisonFS/file"

	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type Filesystem struct {
	filei *file.FileTree
}

func Mount(filei *file.FileTree, mountpoint string) {
	vifs := newVisonFS(filei)
	var opts pathfs.PathNodeFsOptions
	op := nodefs.NewOptions()
	opts.ClientInodes = false
	opts.Debug = false
	fs := pathfs.NewPathNodeFs(vifs, &opts)
	state, _, err := nodefs.MountRoot(mountpoint, fs.Root(), op)
	if err != nil {
		panic(err)
	}
	state.Serve()

}
