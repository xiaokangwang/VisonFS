package cache

import (
	"os"
	"path/filepath"
	"strings"
)

func isExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
func IsExist(path string) bool {
	return isExist(path)
}
func IsDirty(path string) bool {
	return isExist(path + ".dirty")
}
func SetDirty(path string) {
	f, _ := os.Create(path + ".dirty")
	f.Close()
}
func RemoveDirty(path string) {
	os.Remove(path + ".dirty")
}
func Purge(path string) {
	filepath.Walk(path+"/blob/", func(pathi string, f os.FileInfo, err error) error {
		if !IsDirty(pathi) && !strings.HasSuffix(pathi, ".dirty") {
			if f.IsDir() {
				return nil
			}
			err := os.Remove(pathi)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
}
func FindDrity(path string) []string {
	var uploading []string
	filepath.Walk(path+"/blob/", func(pathi string, f os.FileInfo, err error) error {
		if strings.HasSuffix(pathi, ".dirty") {
			s := pathi[:len(pathi)-len(".dirty")]
			l := s[len(path)+1:]
			uploading = append(uploading, l)
		}
		return nil
	})
	return uploading
}
