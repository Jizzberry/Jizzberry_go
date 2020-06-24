package scan

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

const component = "Scan"

type Scan struct {
}

func worker(paths []string, ctx context.Context, progress *int) {
	wg := sync.WaitGroup{}

	maxGoroutines := runtime.NumCPU() / 2
	maxController := make(chan struct{}, maxGoroutines)

	files := make([]string, 0)

	for _, item := range paths {
		tmp, err := tasks_handler.GetAllFiles(item)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		files = append(files, tmp...)
	}
	progressMutex := sync.Mutex{}
	tmp := make(chan int, len(files))
	*progress = 0

	filesModel := files2.Initialize()
	defer filesModel.Close()

	for _, f := range files {
		maxController <- struct{}{}
		wg.Add(1)
		go func(f string, maxController chan struct{}) {
			select {
			case <-ctx.Done():
				wg.Done()
				tmp <- 1
				progressMutex.Lock()
				*progress = int(float32(len(tmp)) / float32(len(files)) * 100)
				progressMutex.Unlock()
				<-maxController
				return
			default:
				file := files2.Files{}

				_, exists := filesModel.IsExists(f)

				if !exists {
					info, _ := os.Stat(f)
					ext := filepath.Ext(f)
					file.FileName = strings.ReplaceAll(info.Name(), ext, "")
					file.Length = ffmpeg.GetLength(f)
					file.FileSize = ByteCountDecimal(info.Size())
					file.DateCreated = time.Unix(info.ModTime().Unix(), 0).Format("01-02-06")
					file.FilePath = f
					file.Tags = ""
					actorsData := tasks_handler.MatchActorToTitle(info.Name())
					var joinActors string
					for i, a := range actorsData {
						scraped := scrapers.ScrapeActor(a)
						actor_details.Initialize().Create(*scraped)
						joinString(&joinActors, a.Name, i != len(actorsData)-1)
					}
					file.Actors = joinActors

					var joinStudios string
					studiosData := tasks_handler.MatchStudioToTitle(info.Name())
					for i, s := range studiosData {
						joinString(&joinStudios, s.Name, i != len(actorsData)-1)
					}

					file.Studios = joinStudios

					genId := filesModel.Create(file)

					ffmpeg.GenerateThumbnail(genId, f, lenInSec(file.Length)/2)

					helpers.LogInfo(fmt.Sprintf("scanned %s successfully", f), component)
				} else {
					helpers.LogInfo(fmt.Sprintf("skipped %s", f), component)
				}
				wg.Done()
				tmp <- 1
				progressMutex.Lock()
				*progress = int(float32(len(tmp)) / float32(len(files)) * 100)
				progressMutex.Unlock()
				<-maxController
			}
		}(f, maxController)
	}
	wg.Wait()
	close(tmp)
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
