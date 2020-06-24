package tasks

import "os"

func makeLink(src string, target string) error {
	return os.Symlink(src, target)
}
