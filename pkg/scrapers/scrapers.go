package scrapers

import (
	"context"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers/factory"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers/pornhub"
	"sync"
)

var actorScrapers = make([]factory.ActorsImpl, 0)
var videoScrapers = make([]factory.VideosImpl, 0)
var studioScrapers = make([]factory.StudiosImpl, 0)

func RegisterScrapers() {
	actorScrapers = append(actorScrapers, pornhub.Pornhub{})
	videoScrapers = append(videoScrapers, pornhub.Pornhub{})
	studioScrapers = append(studioScrapers, pornhub.Pornhub{})
}

func ScrapeActor(sceneId int64, actors actor.Actor) *actor_details.ActorDetails {
	detailsModel := actor_details.Initialize()

	if !detailsModel.IsExists(actors.GeneratedID) {
		for _, i := range actorScrapers {
			if i.GetWebsite() == actors.Website {
				details := i.ScrapeActor(actors.Name)
				i.ScrapeImage(actors.Name, actors.GeneratedID)

				//Manually set name just in case of connection error
				details.Name = actors.Name
				details.SceneId = sceneId
				details.ActorId = actors.GeneratedID
				return &details
			}
		}
	}
	if actors.GeneratedID > 0 {
		details := detailsModel.Get(actor_details.ActorDetails{ActorId: actors.GeneratedID})
		if len(details) > 0 {
			details[0].SceneId = sceneId
			return &details[0]
		}
		return &actor_details.ActorDetails{}
	} else {
		details := actor_details.ActorDetails{}
		details.Name = actors.Name
		details.SceneId = sceneId
		return &details
	}

}

func ScrapeActorList(ctx context.Context, progress *int) {
	tmp := make(chan int, len(actorScrapers))
	progressMutex := sync.Mutex{}
	for _, i := range actorScrapers {
		go func(i factory.ActorsImpl) {
			i.ScrapeActorList(ctx)
			tmp <- 1
			progressMutex.Lock()
			*progress = int(float32(len(tmp)) / float32(len(actorScrapers)) * 100)
			progressMutex.Unlock()
		}(i)
	}
}

func ScrapeStudioList(ctx context.Context, progress *int) {
	tmp := make(chan int, len(studioScrapers))
	progressMutex := sync.Mutex{}
	for _, i := range studioScrapers {
		go func(i factory.StudiosImpl) {
			i.ScrapeStudiosList(ctx)
			tmp <- 1
			progressMutex.Lock()
			*progress = int(float32(len(tmp)) / float32(len(actorScrapers)) * 100)
			progressMutex.Unlock()
		}(i)
	}
}

func ScrapeVideo(url string) factory.VideoDetails {
	for _, i := range videoScrapers {
		if i.ParseUrl(url) {
			return i.ScrapeVideo(url)
		}
	}

	return factory.VideoDetails{}
}

func QueryVideos(query string) map[string][]factory.Videos {
	scrapedVideos := make(map[string][]factory.Videos)
	for _, i := range videoScrapers {
		scrapedVideos[i.GetWebsite()] = i.QueryVideos(query)
	}
	return scrapedVideos
}
