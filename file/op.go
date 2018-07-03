package file

import "github.com/xiaokangwang/VisonFS/journeldb"
import "github.com/xiaokangwang/VisonFS/transform"

type FileTree struct {
	jd *journeldb.JournelDB
	tf *transform.Transform
}

func (ft *FileTree) Ls(path string)                 {}
func (ft *FileTree) GetFileMeta(path string) string {}
func (ft *FileTree) SetFileMeta(path, meta string)  {}

//Block=16MB
//May Block if file is not ready
func (ft *FileTree) GetFileBlock(path string, blockid int) []byte {
}

//May Block if writethrough is true
func (ft *FileTree) SetFileBlock(path string, blockid int, content []byte, writethrough bool) {}

func (ft *FileTree) Mkdir(path, ele string) {}
func (ft *FileTree) Rm(path, ele string)    {}
