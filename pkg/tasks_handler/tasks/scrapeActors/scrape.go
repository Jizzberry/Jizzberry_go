package scrapeActors

import (
	"context"
	"github.com/Jizzberry/Jizzberry_go/pkg/scrapers"
)

type ScrapeActors struct {
}

func (s ScrapeActors) Start() (*context.CancelFunc, *int) {
	var progress int
	ctx, cancel := context.WithCancel(context.Background())
	scrapers.ScrapeActorList(ctx, &progress)
	return &cancel, &progress
}
