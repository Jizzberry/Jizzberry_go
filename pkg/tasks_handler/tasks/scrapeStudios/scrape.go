package scrapeStudios

import (
	"context"
	"github.com/Jizzberry/Jizzberry_go/pkg/scrapers"
)

type ScrapeStudios struct {
}

func (s ScrapeStudios) Start() (*context.CancelFunc, *int) {
	var progress int
	ctx, cancel := context.WithCancel(context.Background())
	scrapers.ScrapeStudioList(ctx, &progress)
	return &cancel, &progress
}
