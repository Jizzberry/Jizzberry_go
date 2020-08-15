package scrapers

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/gocolly/colly/v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func getScrapeImage(i int, actor actor.Actor) (path string) {
	data := helpers.SafeMapCast(helpers.SafeSelectFromMap(ParseYaml(scrapers[i].path), helpers.ScraperImage))
	if data != nil {
		url := helpers.SafeCastString(data[helpers.YamlURL])
		if url != "" {
			var link string
			c := getColly(nil, func(e *colly.HTMLElement) {
				getDataAndScrape(data, helpers.ImageLink, e, &link, func(string) bool { return true })
				path = helpers.GetThumbnailPath()
				downloadImage(link, filepath.Join(helpers.ThumbnailPath, path))
			})

			err := c.Visit(parseUrl(url, actor.UrlID))
			if err != nil {
				helpers.LogError(err.Error())
			}
			c.Wait()
		} else {
			helpers.LogError("Couldn't find url")
		}
	} else {
		helpers.LogError("Couldn't parse data")
	}
	return
}

// Downloads image from link and returns path of out file
func downloadImage(link string, outPath string) {
	if link != "" {
		file, err := os.Create(outPath)
		if err != nil {
			helpers.LogError(err.Error())
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
			helpers.LogError(err.Error())
			return
		}
		defer resp.Body.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			helpers.LogError(err.Error())
		}
	}
}
