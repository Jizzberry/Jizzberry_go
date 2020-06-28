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
				getDataAndScrape(data, helpers.VideoTitle, e, &details.Name, false, func(string) bool { return true })
				getDataAndScrape(data, helpers.VideoActors, e, &details.Actors, true, func(string) bool { return true })
				getDataAndScrape(data, helpers.VideoTags, e, &details.Tags, true, func(string) bool { return true })
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
	if data != nil {
		if url != "" {
			names := make([]string, 0)
			links := make([]string, 0)

			c := getColly(func(e *colly.HTMLElement) {
				getDataAndScrape(data, helpers.VideosName, e, &names, true, func(string) bool { return true })
				getDataAndScrape(data, helpers.VideosLink, e, &links, true, func(string) bool { return true })
				videos = compileResults(names, links, website)
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
