package utils

import "os"

func FileExits(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
