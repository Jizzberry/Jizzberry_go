package rename

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/config"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers/factory"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks"
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

		for _, actor := range tasks.MatchName(a.Name).Actors {
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
	actorDetails.Delete(sceneId)
	for _, a := range actors {
		data := tasks.MatchActorExact(a)
		for _, a := range data.Actors {
			scraped := scrapers.ScrapeActor(sceneId, a)
			actor_details.Initialize().Create(scraped)
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

func (r Rename) Start(sceneId int64, title string, actors []string) {

	title = FormatTitle(title, sceneId, actors)

	folders := getFolder(sceneId, actors, title)
	err := makeFolders(folders)
	if err != nil {
		fmt.Println(err)
		return
	}

	originalPath := files.Initialize().Get(sceneId).FilePath
	ext := filepath.Ext(originalPath)
	for i := range folders {

		if !isFileExists(folders[i]) {
			_ = fmt.Errorf("one of the file already exists")
			return
		}

		var err error
		folders[i], err = filepath.Abs(filepath.FromSlash(folders[i] + "/" + title + ext))
		if err != nil {
			fmt.Println(err)
		}
	}

	err = moveFile(originalPath, folders[0])
	if err != nil {
		fmt.Println(err)
		return
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

func FormatTitle(title string, sceneId int64, actors []string) string {
	formatter := config.GetFileRenameFormatter()
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
			tags := files.Initialize().GetTags(sceneId)
			formatter = strings.ReplaceAll(formatter, m, strings.Join(tags, ", "))
		}
	}
	return formatter
}

func getBasePath(sceneId int64) string {
	scenePath := files.Initialize().Get(sceneId).FilePath
	videoPaths := config.GetVideoPaths()

	for _, p := range videoPaths {
		if strings.Contains(scenePath, p) {
			return p
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
			tags := files.Initialize().GetTags(sceneId)
			for _, t := range tags {
				finalFolders = append(finalFolders, filepath.FromSlash(basePath+"/"+t))
			}
		}
	} else {
		for _, m := range matches {
			if strings.ToLower(m) == "{{actors}}" {
				strings.ReplaceAll(formatter, m, strings.Join(actors, ", "))

			} else if strings.ToLower(m) == "{{title}}" {
				strings.ReplaceAll(formatter, m, title)

			} else if strings.ToLower(m) == "{{tags}}" {
				strings.ReplaceAll(formatter, m, strings.Join(files.Initialize().GetTags(sceneId), ", "))

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
