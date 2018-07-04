package transfer

import (
	"bytes"
	"io"
	"os"

	"github.com/xiaokangwang/VisonFS/file"
)

type Transfer struct {
	File                string
	RFile               string
	LastTransferedBlock int64
	Upload              bool
	filei               *file.FileTree
	Size                int64
}

const Blocksize = 1024 * 1024 * 16

func (t *Transfer) UploadMeta() {
	t.filei.SetSize(t.RFile, t.Size)
}
func (t *Transfer) BlockSum() int64 {
	return (t.Size / Blocksize) + 1
}

func (t *Transfer) ProcessBlock() {
	if !t.HasNext() {
		panic(nil)
	}
	if t.LastBlock() {
		if t.Upload {
			t.UploadMeta()
		} else {
			//none
		}
		t.LastTransferedBlock++
	}
	if t.Upload {
		t.progressUpload()
	} else {
		t.progressDownload()
	}
}
func (t *Transfer) HasNext() bool {
	return t.BlockSum() <= t.LastTransferedBlock
}
func (t *Transfer) LastBlock() bool {
	return t.BlockSum() == t.LastTransferedBlock
}
func (t *Transfer) progressUpload() {
	//Calc next block position
	loc := t.LastTransferedBlock * Blocksize
	lfile, err := os.Open(t.File)
	if err != nil {
		panic(err)
	}
	lfile.Seek(loc, 0)
	r := io.LimitReader(lfile, Blocksize)
	buf := make([]byte, Blocksize)
	io.ReadFull(r, buf)
	lfile.Close()
	t.filei.SetFileBlock(t.RFile, int(t.LastTransferedBlock+1), buf, true)
	t.LastTransferedBlock++
}
func (t *Transfer) progressDownload() {
	block := t.filei.GetFileBlock(t.RFile, int(t.LastTransferedBlock+1))
	lfile, err := os.Open(t.File)
	if err != nil {
		panic(err)
	}
	loc := t.LastTransferedBlock * Blocksize
	lfile.Seek(loc, 0)
	io.Copy(lfile, bytes.NewReader(block))
	lfile.Close()
}
func NewTask(File string,
	RFile string,
	Upload bool, filei *file.FileTree) *Transfer {
	info, err := os.Stat(File)
	if err != nil {
		panic(err)
	}
	return &Transfer{RFile: RFile, File: File, Upload: Upload, filei: filei, Size: info.Size()}
}
func (t *Transfer) PushFileInstance(filei *file.FileTree) {
	t.filei = filei
}
