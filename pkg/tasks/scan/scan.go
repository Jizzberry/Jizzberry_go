package scan

import (
	"context"
	"github.com/Jizzberry/Jizzberry-go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
	files2 "github.com/Jizzberry/Jizzberry-go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Scan struct {
}

func worker(paths []string, ctx context.Context) {
	wg := sync.WaitGroup{}

	maxGoroutines := runtime.NumCPU() / 2
	maxController := make(chan struct{}, maxGoroutines)
	for _, item := range paths {
		files := getAllFiles(item)

		for _, f := range files {
			maxController <- struct{}{}
			wg.Add(1)
			go func(f string, guard chan struct{}) {

				select {
				case <-ctx.Done():
					return
				default:
					filesModel := files2.Initialize()
					files := files2.Files{}

					_, exists := filesModel.IsExists(f)

					if !exists {
						info, _ := os.Stat(f)
						ext := filepath.Ext(f)
						files.FileName = strings.ReplaceAll(info.Name(), ext, "")
						files.Length = ffmpeg.GetLength(f)
						files.FileSize = strconv.FormatInt(info.Size(), 10)
						files.DateCreated = strconv.FormatInt(info.ModTime().Unix(), 10)
						files.FilePath = f
						files.Tags = ""
						genId := filesModel.Create(files)
						ffmpeg.GenerateThumbnail(genId, f, 30)

						data := tasks.MatchName(info.Name())

						for _, a := range data.Actors {
							actor_details.Initialize().Create(*scrapers.ScrapeActor(genId, a))
						}
					}
					wg.Done()
					<-guard
				}
			}(f, maxController)
		}
	}
}

func (s Scan) Start(paths []string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go worker(paths, ctx)
	return ctx, cancel
}

func getAllFiles(path string) []string {
	fileList := make([]string, 0)
	filepath.Walk(path, func(filePath string, f os.FileInfo, err error) error {
		if f.IsDir() == false && isValidExt(filepath.Ext(filePath)) == true {
			fileList = append(fileList, filePath)
		}
		return nil
	})
	return fileList
}

func isValidExt(ext string) bool {
	switch ext {
	case
		".mp4",
		".avi",
		".webm",
		".flv":
		return true
	}
	return false
}
