package utils

import "os"

func IsDirExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		return fileInfo.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}
