package utils

import (
	"os"
	"path/filepath"
)

func IsFileExist(fpath string) bool {
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		return false
	}
	return true
}

func RecursiveFindFiles(root string, pattern string) ([]string, error) {
	files := make([]string, 0)
	findfile := func(path string, f os.FileInfo, err error) (inner error) {
		if err != nil {
			return
		}
		if f.IsDir() {
			return
		} else if match, innerr := filepath.Match(pattern, f.Name()); innerr == nil && match {
			files = append(files, path)
		}
		return
	}
	err := filepath.Walk(root, findfile)
	if len(files) == 0 {
		return files, err
	} else {
		return files, err
	}
}
