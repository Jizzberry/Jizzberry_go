package scrapers

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/studios"
	"github.com/gocolly/colly/v2"
	"strconv"
)

func getScrapeStudiosList(i int) {
	yamlData := ParseYaml(scrapers[i].path)
	data := safeMapCast(safeSelectFromMap(yamlData, helpers.ScraperStudioList))
	if data != nil {
		lastPage := safeConvertInt(data[helpers.YamlLastPage])
		url := safeCastString(safeSelectFromMap(safeMapCast(data), helpers.YamlURL))

		if lastPage < 0 || url == "" {
			helpers.LogError("last_page or url not specified", component)
			return
		}

		model := studios.Initialize()
		defer model.Close()

		c := getColly(func(e *colly.HTMLElement) {
			dest := make([]string, 0)
			getDataAndScrape(data, helpers.StudioListName, e, &dest, true, func(string) bool { return true })
			addStudio(model, dest)
		},
		)

		q := getQueue()
		for i := 1; i < lastPage; i++ {
			url := parseUrl(url, strconv.Itoa(i))
			err := q.AddURL(url)
			if err != nil {
				helpers.LogError(err.Error(), component)
			}
		}
		err := q.Run(c)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
		c.Wait()
	} else {
		helpers.LogError("Couldn't parse data", component)
	}
}

func addStudio(model *studios.Model, strs []string) {
	studioSlice := make([]studios.Studio, 0)
	for _, s := range strs {
		studioSlice = append(studioSlice, studios.Studio{
			Name: s,
		})
	}
	model.Create(studioSlice)
}