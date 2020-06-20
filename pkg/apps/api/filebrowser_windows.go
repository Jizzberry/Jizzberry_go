package api

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"path/filepath"
	"strings"
	"syscall"
)

func getDrives() (drives []string) {

	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	getLogicalDrivesHandle, _ := syscall.GetProcAddress(kernel32, "GetLogicalDrives")

	if ret, _, callErr := syscall.Syscall(getLogicalDrivesHandle, 0, 0, 0, 0); callErr != 0 {
		helpers.LogError(callErr.Error(), component)
	} else {
		drives = bitsToDrives(uint32(ret))
	}

	return
}

func bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i]+":\\")
		}
		bitMap >>= 1
	}
	return
}

func GetDirectory(path string) (b []browse) {
	if len(path) < 1 || path == "/" {
		drives := getDrives()
		for _, d := range drives {
			b = append(b, browse{
				Name: d,
				Path: d,
			})
		}
	} else {
		b = append(b, browse{
			Name: "..",
			Path: func() string {
				if len(clean(strings.Split(path, "\\"))) == 1 {
					return ""
				} else {
					return filepath.Join(path, "..")
				}
			}(),
		})

		allFiles, err := getAllFolders(filepath.FromSlash(path))
		if err != nil {
			return
		}

		for _, f := range allFiles {
			split := strings.Split(f, "\\")
			b = append(b, browse{
				Name: split[len(split)-1],
				Path: f,
			})
		}
	}
	return
}

func clean(a []string) (clean []string) {
	for _, str := range a {
		if str != "" {
			clean = append(clean, str)
		}
	}
	return
}
