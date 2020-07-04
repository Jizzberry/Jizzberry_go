package tasks

import (
	"context"
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor_details"
	files2 "github.com/Jizzberry/Jizzberry_go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry_go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry_go/pkg/tasks_handler"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const component = "Tasks"

type Scan struct {
}

func worker(paths []string, ctx context.Context, progress *int) {
	if paths != nil {
		filesModel := files2.Initialize()
		defer filesModel.Close()

		files := make([]string, 0)

		for _, item := range paths {
			tmp, err := tasks_handler.GetAllFiles(item)
			if err != nil {
				helpers.LogError(err.Error(), component)
			}
			files = append(files, tmp...)
		}

		if len(files) != 0 {
			wg := sync.WaitGroup{}
			maxGoroutines := runtime.NumCPU() / 2
			maxController := make(chan struct{}, maxGoroutines)

			progressMutex := sync.Mutex{}
			tmp := make(chan int, len(files))
			*progress = 0

			for _, f := range files {
				maxController <- struct{}{}
				wg.Add(1)

				go func(f string) {
					select {
					case <-ctx.Done():
						updateProgress(progress, tmp, len(files), &progressMutex)
						wg.Done()
						<-maxController
						return
					default:
						if exists := filesModel.IsExists(f); !exists {
							info, err := os.Stat(f)
							if err != nil {
								helpers.LogError(err.Error(), component)
							}
							thumbnailPath := helpers.GetThumbnailPath()
							file := createFile(f, info, filepath.Ext(f), thumbnailPath)
							filesModel.Create(file)
							ffmpeg.GenerateThumbnail(f, lenInSec(file.Length)/2, thumbnailPath)
							helpers.LogInfo(fmt.Sprintf("scanned %s successfully", f), component)
						} else {
							helpers.LogInfo(fmt.Sprintf("skipped %s", f), component)
						}
						updateProgress(progress, tmp, len(files), &progressMutex)
						wg.Done()
						<-maxController
					}
				}(f)
			}
			wg.Wait()
			close(tmp)
		}
		cleanNoneExistent(filesModel)
		*progress = 100
	}
}

func (s Scan) Start(paths []string) (*context.CancelFunc, *int) {
	var progress int
	ctx, cancel := context.WithCancel(context.Background())
	go worker(paths, ctx, &progress)
	return &cancel, &progress
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func lenInSec(length string) (secs int64) {
	split := strings.Split(length, ":")
	for i := 0; i < len(split); i++ {
		val, err := strconv.ParseInt(split[len(split)-i-1], 10, 64)
		if err != nil {
			helpers.LogError(err.Error(), component)
			continue
		}
		if i == 3 {
			secs += 24 * val
			continue
		}
		secs += int64(float64(val) * math.Pow(60, float64(i)))
	}
	return
}

func joinString(full *string, part string, seperator bool) {
	*full += part
	if seperator {
		*full += ", "
	}
}

func updateProgress(progress *int, current chan int, total int, mutex *sync.Mutex) {
	mutex.Lock()
	current <- 1
	if total != 0 {
		*progress = int(float32(len(current)) / float32(total) * 100)
	}
	mutex.Unlock()
}

func getActors(title string) (joinActors string) {
	actorsData := tasks_handler.MatchActorToTitle(title)
	model := actor_details.Initialize()
	defer model.Close()
	for i, a := range actorsData {
		scraped := scrapers.ScrapeActor(a)
		model.Create(scraped)
		joinString(&joinActors, a.Name, i != len(actorsData)-1)
	}
	return joinActors
}

func getStudios(title string) (joinStudios string) {
	studiosData := tasks_handler.MatchStudioToTitle(title)
	for i, s := range studiosData {
		joinString(&joinStudios, s.Name, i != len(studiosData)-1)
	}
	return
}

func createFile(filepath string, info os.FileInfo, ext string, thumbnailPath string) files2.Files {
	file := files2.Files{}

	file.FileName = strings.ReplaceAll(info.Name(), ext, "")
	file.Length = ffmpeg.GetLength(filepath)
	file.FileSize = ByteCountDecimal(info.Size())
	file.DateCreated = time.Unix(info.ModTime().Unix(), 0).Format(helpers.DateLayout)
	file.FilePath = filepath
	file.Tags = ""
	file.Actors = getActors(info.Name())
	file.Studios = getStudios(info.Name())
	file.ThumbnailPath = thumbnailPath

	return file
}

func cleanNoneExistent(model *files2.Model) {
	files := model.Get(files2.Files{})
	for _, f := range files {
		if _, err := os.Stat(f.FilePath); os.IsNotExist(err) {
			model.Delete(f)
		}
	}
}
