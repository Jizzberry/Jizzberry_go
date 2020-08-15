package scrapers

import (
	"context"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor_details"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
)

func getScrapeActorsList(i int, ctx context.Context) {
	yamlData := ParseYaml(scrapers[i].path)
	website := helpers.SafeCastString(yamlData[helpers.ScraperWebsite])
	data := helpers.SafeMapCast(helpers.SafeSelectFromMap(yamlData, helpers.ScraperActorList))
	if data != nil {
		lastPage := helpers.SafeConvertInt(data[helpers.YamlLastPage])
		url := helpers.SafeCastString(helpers.SafeSelectFromMap(helpers.SafeMapCast(data), helpers.YamlURL))
		selector := helpers.SafeSelectFromMap(data, helpers.YamlForEach)

		if lastPage < 0 || url == "" {
			helpers.LogError("last_page or url not specified")
			return
		}

		model := actor.Initialize()
		defer model.Close()

		c := getColly(ctx, func(e *colly.HTMLElement) {
			headers := []string{helpers.ActorListName, helpers.ActorListURLID}
			dest := make([][]string, len(headers))

			scrapeList(selector, data, headers, &dest, e, func(s string, i int) bool {
				if i == 0 {
					split := len(strings.FieldsFunc(s, splitter))
					return split <= 3
				}
				return true
			})

			addActors(model, dest[0], dest[1], website)
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

func getScrapeActor(i int, actor actor.Actor) (actorDetails actor_details.ActorDetails) {
	data := helpers.SafeMapCast(helpers.SafeSelectFromMap(ParseYaml(scrapers[i].path), helpers.ScraperActor))
	if data != nil {
		url := helpers.SafeCastString(data[helpers.YamlURL])
		if url != "" {
			headers := []string{helpers.ActorName, helpers.ActorBday, helpers.ActorBplace, helpers.ActorHeight, helpers.ActorWeight}
			destinations := []*string{&actorDetails.Name, &actorDetails.Birthday, &actorDetails.Birthplace, &actorDetails.Height, &actorDetails.Weight}

			c := getColly(nil, func(e *colly.HTMLElement) {
				for i := range headers {
					getDataAndScrape(data, headers[i], e, destinations[i], func(string) bool { return true })
				}
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

// Creates array of Actor and adds to db
func addActors(model *actor.Model, strs []string, links []string, website string) {
	if len(strs) == len(links) {
		actorSlice := make([]actor.Actor, 0)
		for i := range strs {
			actorSlice = appendIfNotExists(actorSlice, actor.Actor{
				Name:    strs[i],
				Website: website,
				UrlID:   links[i],
			})
		}
		model.Create(actorSlice)
	} else {
		helpers.LogError("Length of name and links is different")
	}
}

func splitter(r rune) bool {
	return r == ' ' || r == '-' || r == '_'
}

func appendIfNotExists(slice []actor.Actor, actor2 actor.Actor) []actor.Actor {
	for _, a := range slice {
		if a.Name == actor2.Name {
			return slice
		}
	}
	return append(slice, actor2)
}
