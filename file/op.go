package file

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xiaokangwang/VisonFS/protectedFolder"
	"github.com/xiaokangwang/VisonFS/transform"
)

type FileTree struct {
	tf *transform.Transform
	pf *protectedFolder.DelegatedAccess
}

func (ft *FileTree) Ls(path string) ([]os.FileInfo, error) {
	fl, err := ft.pf.ListFile(path)
	if err != nil {
		return nil, err
	}
	for k := range fl {
		fl[k] = &tranFileinfo{inner: fl[k]}
	}
	return fl, nil
}

type tranFileinfo struct {
	inner    os.FileInfo
	truesize int64
}

func (df *tranFileinfo) IsDir() bool        { return !strings.HasSuffix(df.inner.Name(), ".d") }
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
	return nil
}

//May Block if writethrough is true
func (ft *FileTree) SetFileBlock(path string, blockid int, content []byte, writethrough bool) {}

func (ft *FileTree) Mkdir(path, ele string) {}
func (ft *FileTree) Rm(path, ele string)    {}
func (ft *FileTree) GetSize(path string) int64 {
	f, _ := ft.pf.ReadFile(path + "/size")
	s, _ := strconv.ParseInt(string(f), 10, 64)
	return s
}
func (ft *FileTree) SetSize(path string, size int64) {
	s := strconv.FormatInt(size, 10)
	ft.pf.WriteFile(path+"/size", []byte(s))
}
