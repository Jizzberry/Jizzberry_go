package rename

import (
	"context"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/config"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers/factory"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks"
	"github.com/go-ole/go-ole"
	"github.com/zetamatta/go-windows-shortcut"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Rename struct {
}

type bestMatch struct {
	tagLen int
	video  factory.Videos
}

func calcTaglen(query string, results map[string][]factory.Videos) map[string][]bestMatch {
	taglenMatch := make(map[string][]bestMatch)

	for website, videos := range results {
		videoSlice := make([]bestMatch, 0)
		for _, video := range videos {
			splitVideo := strings.FieldsFunc(video.Name, tasks.Splitter)
			splitTitle := strings.FieldsFunc(query, tasks.Splitter)

			for _, word := range splitTitle {
				for i, match := range splitVideo {
					if strings.ToLower(word) == strings.ToLower(match) {
						splitVideo = append(splitVideo[:i], splitVideo[i+1:]...)
					}
				}
			}

			videoSlice = append(videoSlice, bestMatch{
				tagLen: len(splitVideo),
				video:  video,
			})
		}
		taglenMatch[website] = videoSlice
	}
	return taglenMatch
}

func getBestResult(taglenResults map[string][]bestMatch) map[string]factory.Videos {
	topResults := make(map[string]factory.Videos, 0)

	for key, value := range taglenResults {
		if len(value) > 0 {
			min := value[0]
			for _, video := range value {
				if video.tagLen < min.tagLen {
					min = video
				}
			}
			topResults[key] = min.video
		}
	}
	return topResults
}

func GetRenameResult(query string) map[string]factory.VideoDetails {
	detailMap := make(map[string]factory.VideoDetails)
	result := scrapers.QueryVideos(query)
	for website, a := range getBestResult(calcTaglen(query, result)) {
		details := scrapers.ScrapeVideo(a.Url)

		for _, actor := range tasks.MatchName(a.Name) {
			if !sliceContains(details.Actors, actor.Name) {
				details.Actors = append(details.Actors, actor.Name)
			}
		}

		detailMap[website] = details
	}

	return detailMap
}

func makeShortcut(src string, target string) error {
	err := ole.CoInitialize(0)
	if err != nil {
		return err
	}
	defer ole.CoUninitialize()
	if err := shortcut.Make(src, target+".lnk", ""); err != nil {
		return err
	}
	return nil
}

func updateDb(sceneId int64, newPath string, actors []string) error {
	files.Initialize().Update(files.Files{
		GeneratedID: sceneId,
		FilePath:    newPath,
		FileName:    strings.ReplaceAll(filepath.Base(newPath), filepath.Ext(newPath), ""),
	})

	actorDetails := actor_details.Initialize()
	actorDetails.Delete(actor_details.ActorDetails{SceneId: sceneId})
	for _, a := range actors {
		data := tasks.MatchActorExact(a)
		for _, a := range *data {
			scraped := scrapers.ScrapeActor(sceneId, a)
			actor_details.Initialize().Create(*scraped)
		}
	}
	return nil
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

func (r Rename) Start(sceneId int64, title string, actors []string) context.CancelFunc {

	title = FormatTitle(title, sceneId, actors)

	folders := getFolder(sceneId, actors, title)
	err := makeFolders(folders)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	originalPaths := files.Initialize().Get(files.Files{GeneratedID: sceneId})
	if len(originalPaths) > 0 {
		originalPath := originalPaths[0].FilePath

		ext := filepath.Ext(originalPath)
		for i := range folders {

			var err error
			folders[i], err = filepath.Abs(filepath.FromSlash(folders[i] + "/" + title + ext))
			if err != nil {
				fmt.Println(err)
			}

			if isFileExists(folders[i]) {
				fmt.Println("file already exists")
				return nil
			}
		}

		err = moveFile(originalPath, folders[0])
		if err != nil {
			fmt.Println(err)
			return nil
		}

		for i := 1; i < len(folders); i++ {
			err := makeShortcut(folders[0], folders[i])
			if err != nil {
				fmt.Println(err)
			}
		}

		err = updateDb(sceneId, folders[0], actors)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func FormatTitle(title string, sceneId int64, actors []string) string {
	formatter := config.GetFileRenameFormatter()
	if formatter != "" {
		r, err := regexp.Compile("\\{\\{([A-Za-z0-9_]+)\\}\\}")

		if err != nil {
			fmt.Println(err)
			return title
		}

		matches := r.FindAllString(formatter, -1)

		for _, m := range matches {
			if strings.ToLower(m) == "{{actors}}" {
				formatter = strings.ReplaceAll(formatter, m, strings.Join(actors, ", "))
			} else if strings.ToLower(m) == "{{title_full}}" {
				formatter = strings.ReplaceAll(formatter, m, title)
			} else if strings.ToLower(m) == "{{title}}" {
				tmpTitle := title
				for _, a := range actors {
					re := regexp.MustCompile(`(?i)` + a)
					tmpTitle = re.ReplaceAllString(tmpTitle, "")
				}
				formatter = strings.ReplaceAll(formatter, m, strings.TrimSpace(tmpTitle))
			} else if strings.ToLower(m) == "{{tags}}" {
				file := files.Initialize().Get(files.Files{GeneratedID: sceneId})
				if len(file) > 0 {
					formatter = strings.ReplaceAll(formatter, m, file[0].Tags)
				}
			}
		}
		return formatter
	}
	return title
}

func getBasePath(sceneId int64) string {
	scenePaths := files.Initialize().Get(files.Files{GeneratedID: sceneId})
	if len(scenePaths) > 0 {
		scenePath := scenePaths[0].FilePath
		videoPaths := config.GetVideoPaths()

		for _, p := range videoPaths {
			if strings.Contains(scenePath, p) {
				return p
			}
		}
	}
	return ""
}

func getFolder(sceneId int64, actors []string, title string) []string {
	formatter := config.GetFolderRenameFormatter()
	r, err := regexp.Compile("\\{\\{([A-Za-z0-9_]+)\\}\\}")

	if err != nil {
		fmt.Println(err)
	}

	basePath := getBasePath(sceneId)
	matches := r.FindAllString(formatter, -1)

	finalFolders := make([]string, 0)

	if len(matches) == 1 {
		if strings.ToLower(matches[0]) == "{{actors}}" {
			for _, a := range actors {
				finalFolders = append(finalFolders, filepath.FromSlash(basePath+"/"+a))
			}
			return finalFolders

		} else if strings.ToLower(matches[0]) == "{{actors_oneline}}" {
			finalFolders = append(finalFolders, filepath.FromSlash(basePath+"/"+strings.Join(actors, ", ")))
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
		}
	} else {
		for _, m := range matches {
			if strings.ToLower(m) == "{{actors}}" {
				strings.ReplaceAll(formatter, m, strings.Join(actors, ", "))

			} else if strings.ToLower(m) == "{{title}}" {
				strings.ReplaceAll(formatter, m, title)

			} else if strings.ToLower(m) == "{{tags}}" {
				file := files.Initialize().Get(files.Files{GeneratedID: sceneId})
				if len(file) > 0 {
					strings.ReplaceAll(formatter, m, file[0].Tags)
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

func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func isFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
