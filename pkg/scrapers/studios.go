package scrapers

import (
	"context"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/studios"
	"github.com/gocolly/colly/v2"
	"strconv"
)

func getScrapeStudiosList(i int, ctx context.Context) {
	yamlData := ParseYaml(scrapers[i].path)
	data := helpers.SafeMapCast(helpers.SafeSelectFromMap(yamlData, helpers.ScraperStudioList))
	if data != nil {
		lastPage := helpers.SafeConvertInt(data[helpers.YamlLastPage])
		url := helpers.SafeCastString(helpers.SafeSelectFromMap(helpers.SafeMapCast(data), helpers.YamlURL))
		selector := helpers.SafeSelectFromMap(data, helpers.YamlForEach)

		if lastPage < 0 || url == "" {
			helpers.LogError("last_page or url not specified")
			return
		}

		model := studios.Initialize()
		defer model.Close()

		c := getColly(ctx, func(e *colly.HTMLElement) {
			headers := []string{helpers.StudioListName}
			dest := make([][]string, len(headers))
			scrapeList(selector, data, headers, &dest, e, func(s string, i int) bool { return true })
			addStudio(model, dest[0])
		},
		)

		q := getQueue()
		for i := 1; i < lastPage; i++ {
			url := parseUrl(url, strconv.Itoa(i))
			err := q.AddURL(url)
			if err != nil {
				helpers.LogError(err.Error())
			}
		}
		err := q.Run(c)
		if err != nil {
			helpers.LogError(err.Error())
		}
		c.Wait()
	} else {
		helpers.LogError("Couldn't parse data")
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
