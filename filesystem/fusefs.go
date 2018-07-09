package filesystem

import (
	"fmt"
	"time"

	"github.com/hanwen/go-fuse/fuse"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/xiaokangwang/VisonFS/file"
)

func newVisonFS(filei *file.FileTree) pathfs.FileSystem {

	return &visonFS{filei: filei}

}

type visonFS struct {
	filei      *file.FileTree
	openedFile map[string](*visonFile)
}

func (fs *visonFS) SetDebug(debug bool) {}

func (fs *visonFS) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	a := &fuse.Attr{}

	a.Owner = *fuse.CurrentOwner()
	if name == "" {
		a.Mode = fuse.S_IFDIR | 0700
		return a, fuse.OK
	}

	attr, err := fs.filei.Attr(name)
	if err == nil {
		if attr.IsDir() {
			a.Mode = fuse.S_IFDIR | 0700
		} else {
			a.Mode = 0700
			a.Size = uint64(fs.filei.GetSize(name))
		}
		return a, fuse.OK
	}

	return nil, fuse.ENOENT

}

func (fs *visonFS) GetXAttr(name string, attr string, context *fuse.Context) ([]byte, fuse.Status) {

	return nil, fuse.ENOATTR

}

func (fs *visonFS) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {

	return fuse.ENOSYS

}

func (fs *visonFS) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {

	return nil, fuse.ENOSYS

}

func (fs *visonFS) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {

	return fuse.ENOSYS

}

func (fs *visonFS) Readlink(name string, context *fuse.Context) (string, fuse.Status) {

	return "", fuse.ENOSYS

}

func (fs *visonFS) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {

	return fuse.ENOSYS

}

func (fs *visonFS) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	fs.filei.Mkdir(name)
	return fuse.OK

}

func (fs *visonFS) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	fs.filei.Rm(name)
	return fuse.OK

}

func (fs *visonFS) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	fs.filei.Rm(name)
	return fuse.OK

}

func (fs *visonFS) Symlink(value string, linkName string, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visonFS) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visonFS) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visonFS) Chmod(name string, mode uint32, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visonFS) Chown(name string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visonFS) Truncate(name string, offset uint64, context *fuse.Context) (code fuse.Status) {
	_, okerr := fs.filei.Attr(name)
	if okerr != nil {
		f := fs.openfile(name)
		ret := f.Truncate(uint64(offset))
		f.Release()
		return ret
	}

	return fuse.ENOENT

}

func (fs *visonFS) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {

	_, okerr := fs.filei.Attr(name)
	if okerr != nil {
		f := fs.openfile(name)
		return f, fuse.OK
	}
	return nil, fuse.ENOENT

}

func (fs *visonFS) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	res, err := fs.filei.Ls(name)
	if err != nil {
		fmt.Println(err)
		return nil, fuse.EINVAL
	}
	for _, v := range res {
		if v == nil {
			continue
		}
		st := new(fuse.DirEntry)
		st.Name = v.Name()
		if v.IsDir() {
			st.Mode = fuse.S_IFDIR | 0700
		} else {
			st.Mode = 0700
		}
		stream = append(stream, *st)
	}
	println(stream)
	return stream, fuse.OK

}

func (fs *visonFS) OnMount(nodeFs *pathfs.PathNodeFs) {

}

func (fs *visonFS) OnUnmount() {

}

func (fs *visonFS) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visonFS) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {

	f := fs.openfile(name)
	return f, fuse.OK

}

func (fs *visonFS) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visonFS) String() string {

	return "visionFS"

}

func (fs *visonFS) StatFs(name string) *fuse.StatfsOut {

	return nil

}
func (fs *visonFS) openfile(name string) *visonFile {
	if of, ok := fs.openedFile[name]; ok {
		of.opencount++
		return of
	}
	size := fs.filei.GetSize(name)
	if size == -1 {
		size = 0
		fs.filei.SetSize(name, 0)
	}
	file := &visonFile{bufferblock: -1, size: size, path: name, opencount: 1}
	return file
}

type visonFile struct {
	bufferblock int
	buffer      []byte
	bufferdirty bool
	size        int64
	path        string
	opencount   int
	filei       *file.FileTree
	fs          *visonFS
}

// NewDefaultFile returns a File instance that returns ENOSYS for

// every operation.

func (f *visonFile) SetInode(*nodefs.Inode) {

}

func (f *visonFile) InnerFile() nodefs.File {

	return nil

}

func (f *visonFile) String() string {

	return "defaultFile"

}

func (f *visonFile) Read(buf []byte, off int64) (fuse.ReadResult, fuse.Status) {
	thisblock := int(off / Blocksize)
	if thisblock == f.bufferblock {
		//Return data from local buffer

	} else {
		//replace buffer
		f.swapBuffer(thisblock)
	}

	//check off
	if off >= f.size {
		return fuse.ReadResultData(nil), fuse.OK
	}

	//find offset projection with in block
	blockoffset := off % Blocksize
	maxremain := Blocksize - blockoffset
	var readsum = len(buf)
	if readsum > int(maxremain) {
		readsum = int(maxremain)
	}
	sizelimit := f.size - off

	if readsum > int(sizelimit) {
		readsum = int(sizelimit)
	}
	return fuse.ReadResultData(f.buffer[blockoffset : int(blockoffset)+readsum-1]), fuse.OK

}

func (f *visonFile) Write(data []byte, off int64) (uint32, fuse.Status) {
	thisblock := int(off / Blocksize)
	if thisblock == f.bufferblock {
		//Write to local buffer
	} else {
		f.swapBuffer(thisblock)
	}
	f.bufferdirty = true
	if len(f.buffer) != Blocksize {
		old := f.buffer
		f.buffer = make([]byte, Blocksize)
		copy(f.buffer, old)
	}
	//find offset projection with in block
	blockoffset := off % Blocksize
	maxremain := Blocksize - blockoffset
	var copysum = len(data)
	if int64(len(data)) > maxremain {
		copysum = int(maxremain)
	}
	copy(f.buffer[blockoffset:], data[:copysum-1])
	if f.size < int64(off+int64(copysum)) {
		f.size = int64(off + int64(copysum))
	}
	return uint32(copysum), fuse.ENOSYS

}

func (f *visonFile) swapBuffer(block int) {
	//first, swap out old data
	if f.bufferdirty {
		f.filei.SetFileBlock(f.path, f.bufferblock, f.buffer, false)
		if f.size != f.filei.GetSize(f.path) {
			f.filei.SetSize(f.path, f.size)
		}
	}
	f.buffer = f.filei.GetFileBlock(f.path, block)
	if f.buffer == nil {
		f.buffer = make([]byte, Blocksize)
	}
	f.bufferdirty = false
	f.bufferblock = block
}

func (f *visonFile) Flock(flags int) fuse.Status { return fuse.ENOSYS }

func (f *visonFile) Flush() fuse.Status {

	return fuse.OK

}

func (f *visonFile) Release() {
	f.opencount--
	if f.opencount == 0 {
		//TODO SYNC
		f.Fsync(0)
		delete(f.fs.openedFile, f.path)
	}
}

func (f *visonFile) GetAttr(*fuse.Attr) fuse.Status {

	return fuse.ENOSYS

}

func (f *visonFile) Fsync(flags int) (code fuse.Status) {
	if f.bufferdirty {
		f.filei.SetFileBlock(f.path, f.bufferblock, f.buffer, false)
		if f.size != f.filei.GetSize(f.path) {
			f.filei.SetSize(f.path, f.size)
		}
	}
	return fuse.OK

}

func (f *visonFile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {

	return fuse.ENOSYS

}

func (f *visonFile) Truncate(size uint64) fuse.Status {
	f.size = int64(size)
	return fuse.ENOSYS

}

func (f *visonFile) Chown(uid uint32, gid uint32) fuse.Status {

	return fuse.ENOSYS

}

func (f *visonFile) Chmod(perms uint32) fuse.Status {

	return fuse.ENOSYS

}

func (f *visonFile) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) {

	return fuse.ENOSYS

}

const Blocksize = 1024 * 1024 * 16
