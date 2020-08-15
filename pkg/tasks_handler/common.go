package tasks_handler

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/studios"
	tags2 "github.com/Jizzberry/Jizzberry_go/pkg/models/tags"
	"github.com/Jizzberry/Jizzberry_go/pkg/scrapers"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const component = "TasksCommon"

func Splitter(r rune) bool {
	return r == ' ' || r == '.' || r == '-' || r == '_' || r == '[' || r == ']' || r == '(' || r == ')'
}

func MatchStudioToTitle(title string) []studios.Studio {
	split := strings.FieldsFunc(title, Splitter)

	recognisedStudios := make([]studios.Studio, 0)

	studiosModel := studios.Initialize()
	defer studiosModel.Close()

	split = cleanSlice(split)

	allStudios := studiosModel.GetFromTitle(split)

	for _, a := range allStudios {
		regex := RegexpBuilder(a.Name)
		r, err := regexp.Compile(regex)
		if err != nil {
			helpers.LogError(err.Error(), component)
			return recognisedStudios
		}

		matches := r.FindAllString(title, -1)
		if len(matches) > 0 {
			recognisedStudios = append(recognisedStudios, a)
		}
	}
	return recognisedStudios
}

func cleanSlice(slice []string) []string {
	words := make([]string, 0)
	for i := range slice {
		// Avoid articles ig
		if len(slice[i]) > 2 && strings.ToLower(slice[i]) != "the" {
			words = append(words, slice[i])
		}
	}
	return words
}

func MatchActorToTitle(title string) []actor.Actor {
	split := strings.FieldsFunc(title, Splitter)

	recognisedActors := make([]actor.Actor, 0)

	actorsModel := actor.Initialize()
	defer actorsModel.Close()

	split = cleanSlice(split)

	allActors := actorsModel.GetFromTitle(split)

	for _, a := range allActors {
		regex := RegexpBuilder(a.Name)
		r, err := regexp.Compile(regex)
		if err != nil {
			helpers.LogError(err.Error(), component)
			return recognisedActors
		}

		matches := r.FindAllString(title, -1)
		if len(matches) > 0 {
			recognisedActors = append(recognisedActors, a)
		}
	}

	return recognisedActors
}

func RegexpBuilder(name string) string {
	replacer := strings.NewReplacer(" ", "\\s*", "-", "\\s*", "_", "\\s*")
	regex := replacer.Replace(name)
	regex = `(?i)` + regex
	return regex
}

func MatchActorExact(name string) *[]actor.Actor {
	model := actor.Initialize()
	defer model.Close()

	actors := make([]actor.Actor, 0)
	actors = append(actors, model.Get(actor.Actor{Name: name})...)

	return &actors
}

func UpdateDetails(sceneId int64, title string, date string, actors []string, tags []string, studios []string) {
	modelFiles := files.Initialize()
	defer modelFiles.Close()

	modelTags := tags2.Initialize()
	defer modelTags.Close()

	for _, t := range tags {
		if len(modelTags.Get(tags2.Tag{Name: t})) < 1 {
			modelTags.Create(tags2.Tag{Name: t})
		}
	}

	file := modelFiles.Get(files.Files{GeneratedID: sceneId})
	if len(file) > 0 {
		if title != "" {
			file[0].FileName = title
		}

		if date != "" {
			file[0].DateCreated = date
		}

		file[0].Actors = strings.Join(actors, ", ")
		file[0].Tags = strings.Join(tags, ", ")
		file[0].Studios = strings.Join(studios, ", ")
		modelFiles.Update(file[0])

		modelActorD := actor_details.Initialize()
		defer modelActorD.Close()

		for _, a := range actors {
			data := MatchActorExact(a)
			for _, a := range *data {
				scraped := scrapers.ScrapeActor(a)
				modelActorD.Create(scraped)
			}
		}
	}
}

func FormatTitle(title string, sceneId int64) string {
	formatter := helpers.GetConfig().FileRenameFormatter
	if formatter != "" {
		r, err := regexp.Compile("\\{\\{([A-Za-z0-9_]+)\\}\\}")

		if err != nil {
			helpers.LogError(err.Error(), component+" - FormatTitle")
			return title
		}

		matches := r.FindAllString(formatter, -1)
		if len(matches) < 1 {
			return title
		}

		model := files.Initialize()
		defer model.Close()

		file := model.Get(files.Files{GeneratedID: sceneId})
		if len(file) < 1 {
			return title
		}

		for _, m := range matches {
			if strings.ToLower(m) == "{{actors}}" {
				actors := file[0].Actors
				formatter = strings.ReplaceAll(formatter, m, actors)
			} else if strings.ToLower(m) == "{{title_full}}" {
				formatter = strings.ReplaceAll(formatter, m, title)
			} else if strings.ToLower(m) == "{{title}}" {
				tmpTitle := title
				actors := file[0].Actors
				for _, a := range strings.Split(actors, ", ") {
					re := regexp.MustCompile(`(?i)` + a)
					tmpTitle = re.ReplaceAllString(tmpTitle, "")
				}
				formatter = strings.ReplaceAll(formatter, m, strings.TrimSpace(tmpTitle))
			} else if strings.ToLower(m) == "{{tags}}" {
				formatter = strings.ReplaceAll(formatter, m, file[0].Tags)
			} else if strings.ToLower(m) == "{{studios}}" {
				formatter = strings.ReplaceAll(formatter, m, file[0].Studios)
			}
		}
		return formatter
	}
	return title
}

func GetAllFiles(path string) ([]string, error) {
	fileList := make([]string, 0)
	err := filepath.Walk(path, func(filePath string, f os.FileInfo, err error) error {
		if f.IsDir() == false && isValidExt(filepath.Ext(filePath)) == true {
			fileList = append(fileList, filePath)
		}
		return nil
	})

	if err != nil {
		return fileList, err
	}
	return fileList, nil
}

type bestMatch struct {
	tagLen int
	video  scrapers.Videos
}

func calcTaglen(query string, results []scrapers.Videos) (taglenMatch []bestMatch) {
	for _, videos := range results {
		splitVideo := strings.FieldsFunc(videos.Name, Splitter)
		splitTitle := strings.FieldsFunc(query, Splitter)

		for _, word := range splitTitle {
			for i, match := range splitVideo {
				if strings.ToLower(word) == strings.ToLower(match) {
					splitVideo = append(splitVideo[:i], splitVideo[i+1:]...)
				}
			}
		}

		taglenMatch = append(taglenMatch, bestMatch{
			tagLen: len(splitVideo),
			video:  videos,
		})
	}
	return taglenMatch
}

func splitByWebsite(taglenResults []bestMatch) map[string][]bestMatch {
	byWebsite := make(map[string][]bestMatch)
	for _, t := range taglenResults {
		byWebsite[t.video.Website] = append(byWebsite[t.video.Website], t)
	}
	return byWebsite
}

func getBestResult(taglenResults []bestMatch) []scrapers.Videos {
	topResults := make([]scrapers.Videos, 0)

	websiteMap := splitByWebsite(taglenResults)

	for _, value := range websiteMap {
		min := value[0]
		for _, v := range value {
			if v.tagLen < min.tagLen {
				min = v
			}
		}
		topResults = append(topResults, min.video)
	}
	return topResults
}

func GetQueryResult(query string) []scrapers.VideoDetails {
	detailMap := make([]scrapers.VideoDetails, 0)
	result := scrapers.QueryVideos(query)
	for _, a := range getBestResult(calcTaglen(query, result)) {
		details := scrapers.ScrapeVideo(a.Url)

		for _, a := range MatchActorToTitle(a.Name) {
			if !sliceContains(details.Actors, a.Name) {
				details.Actors = append(details.Actors, a.Name)
			}
		}

		detailMap = append(detailMap, details)
	}
	return detailMap
}

func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func isValidExt(ext string) bool {
	switch ext {
	case
		".mp4",
		".avi",
		".mkv",
		".webm",
		".flv":
		return true
	}
	return false
}
