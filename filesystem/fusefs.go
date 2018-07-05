package filesystem

import (
	"fmt"
	"time"

	"github.com/hanwen/go-fuse/fuse"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/xiaokangwang/VisonFS/file"
)

func NewVisionFS() pathfs.FileSystem {

	return (*visionFS)(nil)

}

type visionFS struct {
	filei *file.FileTree
}

func (fs *visionFS) SetDebug(debug bool) {}

func (fs *visionFS) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	a := &fuse.Attr{}

	a.Owner = *fuse.CurrentOwner()
	if name == "" {
		a.Mode = fuse.S_IFDIR | 0700
		return a, fuse.OK
	}

	attr, err := fs.filei.Attr(name)
	if err != nil {
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

func (fs *visionFS) GetXAttr(name string, attr string, context *fuse.Context) ([]byte, fuse.Status) {

	return nil, fuse.ENOATTR

}

func (fs *visionFS) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {

	return fuse.ENOSYS

}

func (fs *visionFS) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {

	return nil, fuse.ENOSYS

}

func (fs *visionFS) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {

	return fuse.ENOSYS

}

func (fs *visionFS) Readlink(name string, context *fuse.Context) (string, fuse.Status) {

	return "", fuse.ENOSYS

}

func (fs *visionFS) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {

	return fuse.ENOSYS

}

func (fs *visionFS) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {

	return fuse.ENOSYS

}

func (fs *visionFS) Unlink(name string, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visionFS) Rmdir(name string, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visionFS) Symlink(value string, linkName string, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visionFS) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visionFS) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visionFS) Chmod(name string, mode uint32, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visionFS) Chown(name string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visionFS) Truncate(name string, offset uint64, context *fuse.Context) (code fuse.Status) {
	_, okerr := fs.filei.Attr(name)
	if okerr != nil {
		fs.filei.SetSize(name, int64(offset))
		return fuse.OK
	}

	return fuse.ENOENT

}

func (fs *visionFS) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {

	return nil, fuse.ENOSYS

}

func (fs *visionFS) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, status fuse.Status) {
	res, err := fs.filei.Ls(name)
	if err != nil {
		fmt.Println(err)
		return nil, fuse.EINVAL
	}
	for _, v := range res {
		st := new(fuse.DirEntry)
		st.Name = v.Name()
		if v.IsDir() {
			st.Mode = fuse.S_IFDIR | 0700
		} else {
			st.Mode = 0700
		}
		stream = append(stream, *st)
	}
	return stream, fuse.OK

}

func (fs *visionFS) OnMount(nodeFs *pathfs.PathNodeFs) {

}

func (fs *visionFS) OnUnmount() {

}

func (fs *visionFS) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visionFS) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {

	return nil, fuse.ENOSYS

}

func (fs *visionFS) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {

	return fuse.ENOSYS

}

func (fs *visionFS) String() string {

	return "visionFS"

}

func (fs *visionFS) StatFs(name string) *fuse.StatfsOut {

	return nil

}

type visonFile struct{}

// NewDefaultFile returns a File instance that returns ENOSYS for

// every operation.

func NewDefaultFile() nodefs.File {

	return (*visonFile)(nil)

}

func (f *visonFile) SetInode(*nodefs.Inode) {

}

func (f *visonFile) InnerFile() nodefs.File {

	return nil

}

func (f *visonFile) String() string {

	return "defaultFile"

}

func (f *visonFile) Read(buf []byte, off int64) (fuse.ReadResult, fuse.Status) {

	return nil, fuse.ENOSYS

}

func (f *visonFile) Write(data []byte, off int64) (uint32, fuse.Status) {

	return 0, fuse.ENOSYS

}

func (f *visonFile) Flock(flags int) fuse.Status { return fuse.ENOSYS }

func (f *visonFile) Flush() fuse.Status {

	return fuse.OK

}

func (f *visonFile) Release() {

}

func (f *visonFile) GetAttr(*fuse.Attr) fuse.Status {

	return fuse.ENOSYS

}

func (f *visonFile) Fsync(flags int) (code fuse.Status) {

	return fuse.ENOSYS

}

func (f *visonFile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {

	return fuse.ENOSYS

}

func (f *visonFile) Truncate(size uint64) fuse.Status {

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
