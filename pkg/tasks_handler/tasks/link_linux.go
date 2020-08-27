package tasks

import "os"

func makeLink(src string, target string) (error, string) {
	return os.Symlink(src, target), target
}
