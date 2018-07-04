package cache

import (
	"os"
	"path/filepath"
	"strings"
)

func isExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	}
	return false
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
	filepath.Walk(path, func(pathi string, f os.FileInfo, err error) error {
		if !IsDirty(pathi) && !strings.HasSuffix(pathi, ".dirty") {
			os.Remove(pathi)
		}
		return nil
	})
}
func FindDrity(path string) []string {
	var uploading []string
	filepath.Walk(path, func(pathi string, f os.FileInfo, err error) error {
		if strings.HasSuffix(pathi, ".dirty") {
			uploading = append(uploading, pathi)
		}
		return nil
	})
	return uploading
}
