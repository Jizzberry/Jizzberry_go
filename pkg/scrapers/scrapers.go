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

	if !detailsModel.IsExists(actors.GeneratedID) {
		if exists, index := MatchWebsite(actors.Website); exists {
			if scrapers[index].Actor {
				scraped := getScrapeActor(index, actors.Name)
				scraped.Name = actors.Name
				scraped.ActorId = actors.GeneratedID

				return scraped
			}
		}
	}
	if actors.GeneratedID > 0 {
		details := detailsModel.Get(actor_details.ActorDetails{ActorId: actors.GeneratedID})
		if len(details) > 0 {
			return details[0]
		}
	}
	details := actor_details.ActorDetails{}
	details.Name = actors.Name
	details.ActorId = actors.GeneratedID
	return details
}

func ScrapeActorList(ctx context.Context, progress *int) {
	progressMutex := sync.Mutex{}
	*progress = 1
	for index, s := range scrapers {
		if s.ActorList {
			go func(index int) {
				getScrapeActorsList(index)
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
				getScrapeStudiosList(i)
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
