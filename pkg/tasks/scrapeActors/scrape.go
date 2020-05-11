package scrapeActors

import (
	"context"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers"
)

type ScrapeActors struct {
}

func (s ScrapeActors) Start() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	scrapers.ScrapeActorList(ctx)
	return cancel
}
