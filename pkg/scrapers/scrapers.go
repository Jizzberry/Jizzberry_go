package scrapers

import (
	"context"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor_details"
	"sync"
)

func ScrapeActor(actors actor.Actor) actor_details.ActorDetails {
	detailsModel := actor_details.Initialize()
	defer detailsModel.Close()

	// If actor doesn't already exist in actor_details, scrape them
	if !detailsModel.IsExists(actors.ActorID) {
		if exists, index := MatchWebsite(actors.Website); exists {
			if scrapers[index].Actor {
				scraped := getScrapeActor(index, actors)
				scraped.ThumbnailPath = getScrapeImage(index, actors)
				scraped.Name = actors.Name
				scraped.ActorId = actors.ActorID

				return scraped
			}
		}
	}

	// Avoid scraping if actor already exists in actor_details
	if actors.ActorID > 0 {
		details := detailsModel.Get(actor_details.ActorDetails{ActorId: actors.ActorID})
		if len(details) > 0 {
			return details[0]
		}
	}

	return actor_details.ActorDetails{}
}

func ScrapeActorList(ctx context.Context, progress *int) {
	progressMutex := sync.Mutex{}
	*progress = 1
	for index, s := range scrapers {
		if s.ActorList {
			go func(index int) {
				getScrapeActorsList(index, ctx)
				progressMutex.Lock()
				*progress = int(float32(index) / float32(len(scrapers)) * 100)
				progressMutex.Unlock()
			}(index)
		}
	}
}

func ScrapeStudioList(ctx context.Context, progress *int) {
	progressMutex := sync.Mutex{}
	for i, s := range scrapers {
		if s.StudioList {
			go func(i int) {
				getScrapeStudiosList(i, ctx)
				progressMutex.Lock()
				*progress = int(float32(i) / float32(len(scrapers)) * 100)
				progressMutex.Unlock()
			}(i)
		}
	}
}

func ScrapeVideo(url string) VideoDetails {
	for i, s := range scrapers {
		if s.Video {
			if matchUrlToScraper(url) {
				return getScrapedVideo(i, url)
			}
		}
	}
	return VideoDetails{}
}

func QueryVideos(query string) []Videos {
	scrapedVideos := make([]Videos, 0)
	for i, s := range scrapers {
		if s.QueryVideos {
			scrapedVideos = append(scrapedVideos, getQueryVideo(i, query)...)
		}
	}
	return scrapedVideos
}
