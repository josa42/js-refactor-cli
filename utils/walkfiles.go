package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

var ignoreDirs = []string{"node_modules", ".git"}

func WalkFiles(rootPath string, fn func(path string)) {

	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			_, name := filepath.Split(info.Name())

			for _, iname := range ignoreDirs {
				if name == iname {
					return filepath.SkipDir
				}
			}
			return nil
		}

		fn(path)

		return nil
	})
}
