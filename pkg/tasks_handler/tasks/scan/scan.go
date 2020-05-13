package scan

import (
	"context"
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/ffmpeg"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
	files2 "github.com/Jizzberry/Jizzberry-go/pkg/models/files"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
	"github.com/Jizzberry/Jizzberry-go/pkg/tasks_handler/tasks"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type Scan struct {
}

func worker(paths []string, ctx context.Context, progress *int) {
	wg := sync.WaitGroup{}

	maxGoroutines := runtime.NumCPU() / 2
	maxController := make(chan struct{}, maxGoroutines)

	files := make([]string, 0)

	for _, item := range paths {
		files = append(files, getAllFiles(item)...)
	}
	progressMutex := sync.Mutex{}
	tmp := make(chan int, len(files))
	*progress = 0

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
				fmt.Println("scanning")
				filesModel := files2.Initialize()
				file := files2.Files{}

				_, exists := filesModel.IsExists(f)

				if !exists {
					info, _ := os.Stat(f)
					ext := filepath.Ext(f)
					file.FileName = strings.ReplaceAll(info.Name(), ext, "")
					file.Length = ffmpeg.GetLength(f)
					file.FileSize = strconv.FormatInt(info.Size(), 10)
					file.DateCreated = strconv.FormatInt(info.ModTime().Unix(), 10)
					file.FilePath = f
					file.Tags = ""
					genId := filesModel.Create(file)
					ffmpeg.GenerateThumbnail(genId, f, 30)

					data := tasks.MatchName(info.Name())

					for _, a := range data {
						actor_details.Initialize().Create(*scrapers.ScrapeActor(genId, a))
					}
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
