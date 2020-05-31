package pornhub

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/scrapers/factory"
	"github.com/gocolly/colly/v2"
	"regexp"
	"strings"
)

func (p Pornhub) ScrapeVideo(url string) factory.VideoDetails {
	c := colly.NewCollector(colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"))

	details := factory.VideoDetails{}

	c.OnHTML("body", func(element *colly.HTMLElement) {
		details.Name = element.ChildText(".title-container")

		actors := make([]string, 0)
		element.ForEach(".pornstarsWrapper", func(i int, element *colly.HTMLElement) {
			actors = append(actors, element.ChildAttrs(`a[data-mxptype="Pornstar"]`, "data-mxptext")...)
		})
		details.Actors = actors

		tags := make([]string, 0)
		element.ForEach(".tagsWrapper", func(i int, element *colly.HTMLElement) {
			tags = append(tags, element.ChildTexts(`a`)...)
			tags = tags[:len(tags)-1]
		})

		details.Tags = tags
	})

	err := c.Visit(strings.TrimSpace(url))
	if err != nil {
		helpers.LogError(err.Error(), p.GetWebsite())
	}
	c.Wait()

	return details
}

func (p Pornhub) QueryVideos(query string) []factory.Videos {
	c := colly.NewCollector(colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"))

	videos := make([]factory.Videos, 0)

	c.OnHTML(".nf-videos", func(element *colly.HTMLElement) {
		element.ForEach(".pcVideoListItem", func(i int, element *colly.HTMLElement) {
			video := factory.Videos{}
			video.Name = element.ChildAttr(".title > a", "title")
			video.Url = "https://www.pornhub.com" + element.ChildAttr(".title > a", "href")
			video.Website = p.GetWebsite()
			videos = append(videos, video)
		})
	})

	query = strings.ReplaceAll(strings.ReplaceAll(query, " ", "+"), "-", "+")
	err := c.Visit("https://www.pornhub.com/video/search?search=" + query)
	if err != nil {
		helpers.LogError(err.Error(), p.GetWebsite())
	}

	return videos
}

func (p Pornhub) ParseUrl(url string) bool {
	found, err := regexp.MatchString("pornhub", url)
	if err != nil {
		helpers.LogError(err.Error(), p.GetWebsite()+" Scraper")
	}
	return found
}
