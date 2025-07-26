package utils

import (
	"os"
	"path"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func FileCreatesIfNotExists(filePath string) error {
	if !FileExists(filePath) {
		if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
			return err
		}
		fp, err := os.OpenFile(filePath, os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer fp.Close()
	}
	return nil
}
