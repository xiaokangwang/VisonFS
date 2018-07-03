package file

type FileTree struct {
}

func (ft *FileTree) Ls() {}

//Block=16MB
func (ft *FileTree) GetFileBlock() {}
func (ft *FileTree) SetFileBlock() {}
func (ft *FileTree) Mkdir()        {}
func (ft *FileTree) Rm()           {}
