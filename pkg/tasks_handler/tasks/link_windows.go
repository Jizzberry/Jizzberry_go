package tasks

import (
	"github.com/go-ole/go-ole"
	shortcut "github.com/zetamatta/go-windows-shortcut"
)

func makeLink(src string, target string) (error, string) {
	err := ole.CoInitialize(0)
	if err != nil {
		return err, ""
	}
	defer ole.CoUninitialize()
	if err := shortcut.Make(src, target+".lnk", ""); err != nil {
		return err, ""
	}
	return nil, target + ".lnk"
}
