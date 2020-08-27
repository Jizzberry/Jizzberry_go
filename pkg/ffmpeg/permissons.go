// +build linux

package ffmpeg

import "os"

func setExecutablePerms(filepath string) error {
	return os.Chmod(filepath, 0755)
}
