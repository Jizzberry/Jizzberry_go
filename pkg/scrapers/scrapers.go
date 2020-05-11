package scrapers

import (
	"context"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers/factory"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers/pornhub"
)

var actorScrapers = make([]factory.ActorsImpl, 0)
var videoScrapers = make([]factory.VideosImpl, 0)

func RegisterScrapers() {
	actorScrapers = append(actorScrapers, pornhub.Pornhub{})
	videoScrapers = append(videoScrapers, pornhub.Pornhub{})
}

func ScrapeActor(sceneId int64, actors actor.Actor) *actor_details.ActorDetails {
	detailsModel := actor_details.Initialize()

	if !detailsModel.IsExists(actors.GeneratedID) {
		for _, i := range actorScrapers {
			if i.GetWebsite() == actors.Website {
				details := i.ScrapeActor(actors.Name)

				//Manually set name just in case of connection error
				details.Name = actors.Name
				details.SceneId = sceneId
				details.ActorId = actors.GeneratedID
				return &details
			}
		}
	}
	if actors.GeneratedID > 0 {
		details := detailsModel.Get(actors.GeneratedID)
		details.SceneId = sceneId
		return details
	} else {
		details := actor_details.ActorDetails{}
		details.Name = actors.Name
		details.SceneId = sceneId
		return &details
	}

}

func ScrapeActorList(ctx context.Context) {
	for _, i := range actorScrapers {
		go func(i factory.ActorsImpl) {
			i.ScrapeActorList(ctx)
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
