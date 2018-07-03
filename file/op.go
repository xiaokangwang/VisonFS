package file

import "github.com/xiaokangwang/VisonFS/journeldb"
import "github.com/xiaokangwang/VisonFS/transform"

type FileTree struct {
	jd journeldb.JournelDB
	tf transform.Transform
}

func (ft *FileTree) Ls(path string) {}

//Block=16MB
func (ft *FileTree) GetFileBlock(path string, blockid int) []byte {
}
func (ft *FileTree) SetFileBlock(path string, blockid int, content []byte) {}

func (ft *FileTree) Mkdir(path, ele string) {}
func (ft *FileTree) Rm(path, ele string)    {}
