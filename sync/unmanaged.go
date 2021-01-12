package sync

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Unmanaged(dirs []string, managedDirs []string) []string {

	allDirs := map[string]bool{}
	var unmanagedDirs []string

	for _, dir := range dirs {

		err := filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() && info.Name() == ".git" {
					allDirs[strings.Replace(path, "/.git", "", -1)] = false
					return filepath.SkipDir
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
	}

	for _, managedDir := range managedDirs {
		allDirs[managedDir] = true
	}

	for dir, managed := range allDirs {
		if !managed {
			unmanagedDirs = append(unmanagedDirs, dir)
		}
	}

	return unmanagedDirs
}
