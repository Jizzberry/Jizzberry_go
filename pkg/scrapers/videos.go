package scrapers

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/gocolly/colly/v2"
	"regexp"
)

func getScrapedVideo(i int, url string) (details VideoDetails) {
	yamlData := ParseYaml(scrapers[i].path)
	data := helpers.SafeMapCast(helpers.SafeSelectFromMap(yamlData, helpers.ScraperSingleVideo))
	website := helpers.SafeCastString(helpers.SafeSelectFromMap(yamlData, helpers.ScraperWebsite))

	if data != nil {
		if url != "" {
			c := getColly(nil, func(e *colly.HTMLElement) {
				getDataAndScrape(data, helpers.VideoTitle, e, &details.Name, func(string) bool { return true })
				actors := make([][]string, 1)
				scrapeList(helpers.SafeSelectFromMap(helpers.SafeMapCast(helpers.SafeSelectFromMap(data, helpers.VideoActors)), helpers.YamlForEach), data, []string{helpers.VideoActors}, &actors, e, func(s string, i int) bool { return true })
				tags := make([][]string, 1)
				scrapeList(helpers.SafeSelectFromMap(helpers.SafeMapCast(helpers.SafeSelectFromMap(data, helpers.VideoTags)), helpers.YamlForEach), data, []string{helpers.VideoTags}, &tags, e, func(s string, i int) bool { return true })

				details.Actors = actors[0]
				details.Tags = tags[0]
				details.Website = website
				details.Url = url
			})

			err := c.Visit(url)
			if err != nil {
				helpers.LogError(err.Error(), component)
			}
			c.Wait()
		}
	}
	return
}

func getQueryVideo(i int, query string) (videos []Videos) {
	yamlData := ParseYaml(scrapers[i].path)
	website := helpers.SafeCastString(helpers.SafeSelectFromMap(yamlData, helpers.ScraperWebsite))
	data := helpers.SafeMapCast(helpers.SafeSelectFromMap(yamlData, helpers.ScraperVideos))
	url := parseUrl(helpers.SafeCastString(helpers.SafeSelectFromMap(data, helpers.YamlURL)), query)
	selector := helpers.SafeSelectFromMap(data, helpers.YamlForEach)
	if data != nil {
		if url != "" {
			c := getColly(nil, func(e *colly.HTMLElement) {
				headers := []string{helpers.VideosName, helpers.VideosLink}
				dest := make([][]string, len(headers))
				scrapeList(selector, data, headers, &dest, e, func(s string, i int) bool { return true })
				videos = compileResults(dest[0], dest[1], website)
			})
			err := c.Visit(url)
			if err != nil {
				helpers.LogError(err.Error(), component)
			}
			c.Wait()
		} else {
			helpers.LogError("Couldn't find url", component)
		}
	} else {
		helpers.LogError("Couldn't parse data", component)
	}
	return
}

func matchUrlToScraper(url string) bool {
	for i := range scrapers {
		urlRegex := helpers.SafeCastString(helpers.SafeSelectFromMap(helpers.SafeMapCast(helpers.SafeMapCast(ParseYaml(scrapers[i].path))[helpers.ScraperSingleVideo]), helpers.YamlUrlRegex))
		r, err := regexp.Compile(urlRegex)
		if r != nil {
			if err != nil {
				helpers.LogError(err.Error(), component)
			}
			return r.MatchString(url)
		}
	}
	return false
}

func compileResults(names []string, links []string, website string) (videos []Videos) {
	if len(names) == len(links) {
		for i := range names {
			videos = append(videos, makeVideoStruct(names[i], links[i], website))
		}
	} else {
		helpers.LogError("Length of title and links is different", component)
	}
	return
}

func makeVideoStruct(name string, link string, website string) (videos Videos) {
	videos.Name = name
	videos.Url = link
	videos.Website = website
	return
}
