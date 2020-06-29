package scrapers

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/gocolly/colly/v2"
	"regexp"
)

func getScrapedVideo(i int, url string) (details VideoDetails) {
	data := safeMapCast(safeSelectFromMap(ParseYaml(scrapers[i].path), helpers.ScraperSingleVideo))

	if data != nil {
		if url != "" {
			c := getColly(func(e *colly.HTMLElement) {
				getDataAndScrape(data, helpers.VideoTitle, e, &details.Name, func(string) bool { return true })
				actors := make([][]string, 1)
				scrapeList(safeSelectFromMap(safeMapCast(safeSelectFromMap(data, helpers.VideoActors)), helpers.YamlForEach), data, []string{helpers.VideoActors}, &actors, e, func(s string, i int) bool { return true })
				tags := make([][]string, 1)
				scrapeList(safeSelectFromMap(safeMapCast(safeSelectFromMap(data, helpers.VideoTags)), helpers.YamlForEach), data, []string{helpers.VideoTags}, &tags, e, func(s string, i int) bool { return true })

				details.Actors = actors[0]
				details.Tags = tags[0]
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
	website := safeCastString(safeSelectFromMap(yamlData, helpers.ScraperWebsite))
	data := safeMapCast(safeSelectFromMap(yamlData, helpers.ScraperVideos))
	url := parseUrl(safeCastString(safeSelectFromMap(data, helpers.YamlURL)), query)
	selector := safeSelectFromMap(data, helpers.YamlForEach)
	if data != nil {
		if url != "" {
			c := getColly(func(e *colly.HTMLElement) {
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
		urlRegex := safeCastString(safeSelectFromMap(safeMapCast(safeMapCast(ParseYaml(scrapers[i].path))[helpers.ScraperSingleVideo]), helpers.YamlUrlRegex))
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
