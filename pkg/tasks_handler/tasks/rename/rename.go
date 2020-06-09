package rename

import (
	"context"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const component = "Rename"

type Rename struct {
}

func moveFile(src, target string) error {
	if filepath.VolumeName(src) == filepath.VolumeName(target) {
		err := os.Rename(src, target)
		if err != nil {
			return err
		}
		return nil
	} else {
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		targetFile, err := os.Create(target)
		if err != nil {
			srcFile.Close()
			return err
		}
		defer targetFile.Close()
		_, err = io.Copy(targetFile, srcFile)
		srcFile.Close()
		if err != nil {
			return err
		}
		err = os.Remove(src)
		if err != nil {
			return err
		}
		return nil
	}
}

func makeFolders(folders []string) error {
	for i := range folders {
		err := os.MkdirAll(filepath.FromSlash(folders[i]), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func organize(sceneId int64, progress *int) {
	*progress = 1
	file := files.Initialize().Get(files.Files{GeneratedID: sceneId})
	if len(file) < 1 {
		return
	}

	title := tasks_handler.FormatTitle(file[0].FileName, sceneId)

	folders := getFolder(sceneId, title)
	err := makeFolders(folders)
	if err != nil {
		helpers.LogError(err.Error(), component)
		*progress = 100
		return
	}

	originalPaths := files.Initialize().Get(files.Files{GeneratedID: sceneId})
	if len(originalPaths) > 0 {
		originalPath := originalPaths[0].FilePath

		ext := filepath.Ext(originalPath)
		for i := range folders {

			var err error
			folders[i], err = filepath.Abs(filepath.FromSlash(folders[i] + "/" + title + ext))
			if err != nil {
				helpers.LogError(err.Error(), component)
				*progress = 100
				return
			}

			if isFileExists(folders[i]) {
				helpers.LogError(fmt.Sprintf("file already exists: %s", folders[i]), component)
				*progress = 100
				return
			}
		}

		err = moveFile(originalPath, folders[0])
		if err != nil {
			helpers.LogError(err.Error(), component)
			*progress = 100
			return
		}

		for i := 1; i < len(folders); i++ {
			err := makeLink(folders[0], folders[i])
			if err != nil {
				helpers.LogError(err.Error(), component)
			}
		}
		*progress = 100
	}
}

func (r Rename) Start(sceneId int64) (*context.CancelFunc, *int) {
	var progress int
	go organize(sceneId, &progress)
	return nil, &progress
}

func getBasePath(sceneId int64) string {
	scenePaths := files.Initialize().Get(files.Files{GeneratedID: sceneId})
	if len(scenePaths) > 0 {
		scenePath := scenePaths[0].FilePath
		videoPaths := helpers.GetConfig().Paths

		for _, p := range videoPaths {
			if strings.Contains(scenePath, p) {
				return p
			}
		}
	}
	return ""
}

func getFolder(sceneId int64, title string) []string {
	formatter := helpers.GetConfig().FolderRenameFormatter
	r, err := regexp.Compile("\\{\\{([A-Za-z0-9_]+)\\}\\}")

	if err != nil {
		helpers.LogError(err.Error(), component+" - getFolder")
	}

	basePath := getBasePath(sceneId)
	matches := r.FindAllString(formatter, -1)

	finalFolders := make([]string, 0)

	if len(matches) == 1 {
		if strings.ToLower(matches[0]) == "{{actors}}" {
			actors := files.Initialize().Get(files.Files{GeneratedID: sceneId})[0].Actors
			for _, a := range strings.Split(actors, ", ") {
				finalFolders = append(finalFolders, filepath.FromSlash(basePath+"/"+a))
			}
			return finalFolders

		} else if strings.ToLower(matches[0]) == "{{actors_oneline}}" {
			actors := files.Initialize().Get(files.Files{GeneratedID: sceneId})[0].Actors
			finalFolders = append(finalFolders, filepath.FromSlash(basePath+"/"+actors))
			return finalFolders

		} else if strings.ToLower(matches[0]) == "{{title}}" {
			finalFolders = append(finalFolders, filepath.FromSlash(basePath+"/"+title))
			return finalFolders

		} else if strings.ToLower(matches[0]) == "{{tags}}" {
			file := files.Initialize().Get(files.Files{GeneratedID: sceneId})
			if len(file) > 0 {
				tags := strings.Split(file[0].Tags, ", ")
				for _, t := range tags {
					finalFolders = append(finalFolders, filepath.FromSlash(basePath+"/"+t))
				}
			}
		} else if strings.ToLower(matches[0]) == "{{studios}}" {
			file := files.Initialize().Get(files.Files{GeneratedID: sceneId})
			if len(file) > 0 {
				studios := strings.Split(file[0].Studios, ", ")
				for _, s := range studios {
					finalFolders = append(finalFolders, filepath.FromSlash(basePath+"/"+s))
				}
			}
		}
	} else {
		for _, m := range matches {
			if strings.ToLower(m) == "{{actors}}" {
				actors := files.Initialize().Get(files.Files{GeneratedID: sceneId})[0].Actors
				strings.ReplaceAll(formatter, m, actors)

			} else if strings.ToLower(m) == "{{title}}" {
				strings.ReplaceAll(formatter, m, title)
			} else if strings.ToLower(m) == "{{tags}}" {
				file := files.Initialize().Get(files.Files{GeneratedID: sceneId})
				if len(file) > 0 {
					strings.ReplaceAll(formatter, m, file[0].Tags)
				}
			} else if strings.ToLower(m) == "{{studios}}" {
				file := files.Initialize().Get(files.Files{GeneratedID: sceneId})
				if len(file) > 0 {
					strings.ReplaceAll(formatter, m, file[0].Studios)
				}
			}
		}

		finalFolders = append(finalFolders, filepath.FromSlash(basePath+"/"+formatter))
	}

	if len(finalFolders) == 0 {
		return append(finalFolders, basePath)
	}
	return finalFolders
}

func isFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
