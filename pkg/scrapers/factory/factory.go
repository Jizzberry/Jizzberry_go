package factory

import (
	"context"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
)

type ActorsImpl interface {
	GetWebsite() string
	ScrapeActor(name string) actor_details.ActorDetails
	ScrapeActorList(ctx context.Context)
}

type VideosImpl interface {
	GetWebsite() string
	ScrapeVideo(url string) VideoDetails
	QueryVideos(query string) []Videos
	ParseUrl(url string) bool
}

type VideoDetails struct {
	Name   string
	Actors []string
	Tags   []string
}

type Videos struct {
	Name    string
	Url     string
	Website string
}
