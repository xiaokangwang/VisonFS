package file

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xiaokangwang/VisonFS/gitctl"
	"github.com/xiaokangwang/VisonFS/protectedFolder"
	"github.com/xiaokangwang/VisonFS/sync"
	"github.com/xiaokangwang/VisonFS/transform"
)

type FileTree struct {
	tf      *transform.Transform
	pf      *protectedFolder.DelegatedAccess
	sy      *sync.PendingSync
	gitctli *gitctl.Gitctl
}

func NewFileTree(tf *transform.Transform,
	pf *protectedFolder.DelegatedAccess,
	sy *sync.PendingSync, gitctli *gitctl.Gitctl) *FileTree {
	return &FileTree{tf: tf, pf: pf, sy: sy, gitctli: gitctli}
}

func (ft *FileTree) Ls(path string) ([]os.FileInfo, error) {
	fl, err := ft.pf.ListFile(path)
	if err != nil {
		return nil, err
	}
	for k := range fl {
		if fl[k].Name() != "dir" {
			fl[k] = &tranFileinfo{inner: fl[k]}
			if fl[k].Size() == -1 {
				fl[k] = nil
			}
		} else {
			fl[k] = nil
		}

	}
	return fl, nil
}

type tranFileinfo struct {
	inner    os.FileInfo
	truesize int64
}

func (df *tranFileinfo) IsDir() bool {
	return !strings.HasSuffix(
		df.inner.Name(), ".d")
}
func (df *tranFileinfo) ModTime() time.Time { return df.inner.ModTime() }
func (df *tranFileinfo) Mode() os.FileMode  { return df.inner.Mode() }
func (df *tranFileinfo) Name() string {
	if !df.IsDir() {
		ddn := df.inner.Name()
		return ddn[:len(ddn)-3]
	}
	return df.inner.Name()
}
func (df *tranFileinfo) Size() int64 {
	if !df.IsDir() {
		return df.truesize
	}
	return df.inner.Size()
}
func (df *tranFileinfo) Sys() interface{} { return df.inner.Sys() }

//Block=16MB
//May Block if file is not ready
func (ft *FileTree) GetFileBlock(path string, blockid int) []byte {
	ctx, err := ft.pf.ReadFile(path + ".d/" + strconv.Itoa(blockid))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	b := ft.sy.BlobGet(string(ctx))
	return b
}

//May Block if writethrough is true
func (ft *FileTree) SetFileBlock(path string, blockid int, content []byte, writethrough bool) {
	tracker := ft.sy.BlobUpload(content)
	ft.pf.WriteFile(path+".d/"+strconv.Itoa(blockid), []byte(tracker))
	ft.gitctli.NewVerison()
}

func (ft *FileTree) Mkdir(path string) {
	ft.pf.WriteFile(path+"/dir", []byte("dir"))
	ft.gitctli.NewVerison()
}
func (ft *FileTree) Rm(path string) {
	ft.pf.RemoveFile(path)
	ft.SetSize(path, -1)
	ft.gitctli.NewVerison()
}
func (ft *FileTree) GetSize(path string) int64 {
	f, _ := ft.pf.ReadFile(path + "/size")
	s, _ := strconv.ParseInt(string(f), 10, 64)
	return s
}
func (ft *FileTree) SetSize(path string, size int64) {
	s := strconv.FormatInt(size, 10)
	ft.pf.WriteFile(path+"/size", []byte(s))
	ft.gitctli.NewVerison()
}

func (ft *FileTree) Attr(path string) (os.FileInfo, error) {
	info, err := ft.pf.FileAttr(path)
	if err != nil {
		info, err = ft.pf.FileAttr(path + ".d")
		if err != nil {
			return nil, err
		}
		return &tranFileinfo{inner: info}, nil
	}
	return &tranFileinfo{inner: info}, nil

}
