package scrapers

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor_details"
	"github.com/gocolly/colly/v2"
	"io"
	"net/http"
	"os"
)

func getScrapeImage(i int, actor actor.Actor) (actorDetails actor_details.ActorDetails) {
	data := safeMapCast(safeSelectFromMap(ParseYaml(scrapers[i].path), helpers.ScraperImage))
	if data != nil {
		url := safeCastString(data[helpers.YamlURL])
		if url != "" {
			var link string
			c := getColly(func(e *colly.HTMLElement) {
				getDataAndScrape(data, helpers.ImageLink, e, &link, func(string) bool { return true })
				downloadImage(link, helpers.GetThumbnailPath(actor.GeneratedID, true))
			})

			err := c.Visit(parseUrl(url, actor.Name))
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

func downloadImage(link string, outPath string) {
	if link != "" {
		file, err := os.Create(outPath)
		if err != nil {
			helpers.LogError(err.Error(), component)
			return
		}
		defer file.Close()

		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		resp, err := client.Get(link)
		if err != nil {
			helpers.LogError(err.Error(), component)
			return
		}
		defer resp.Body.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}
}
