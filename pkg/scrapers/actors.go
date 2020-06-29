package scrapers

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor_details"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
)

func getScrapeActorsList(i int) {
	yamlData := ParseYaml(scrapers[i].path)
	website := safeCastString(yamlData[helpers.ScraperWebsite])
	data := safeMapCast(safeSelectFromMap(yamlData, helpers.ScraperActorList))
	if data != nil {
		lastPage := safeConvertInt(data[helpers.YamlLastPage])
		url := safeCastString(safeSelectFromMap(safeMapCast(data), helpers.YamlURL))

		if lastPage < 0 || url == "" {
			helpers.LogError("last_page or url not specified", component)
			return
		}

		model := actor.Initialize()
		defer model.Close()

		c := getColly(func(e *colly.HTMLElement) {
			dest := make([]string, 0)
			getDataAndScrape(data, helpers.ActorListName, e, &dest, true, func(s string) bool {
				split := len(strings.FieldsFunc(s, splitter))
				return split <= 3 && split > 1
			})
			addActors(model, dest, website)
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

func getScrapeActor(i int, actor actor.Actor) (actorDetails actor_details.ActorDetails) {
	data := safeMapCast(safeSelectFromMap(ParseYaml(scrapers[i].path), helpers.ScraperActor))
	if data != nil {
		url := safeCastString(data[helpers.YamlURL])
		if url != "" {
			headers := []string{helpers.ActorName, helpers.ActorBday, helpers.ActorBplace, helpers.ActorHeight, helpers.ActorWeight}
			destinations := []interface{}{&actorDetails.Name, &actorDetails.Birthday, &actorDetails.Birthplace, &actorDetails.Height, &actorDetails.Weight}

			c := getColly(func(e *colly.HTMLElement) {
				for i := range headers {
					getDataAndScrape(data, headers[i], e, destinations[i], false, func(string) bool { return true })
				}
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

func addActors(model *actor.Model, strs []string, website string) {
	actorSlice := make([]actor.Actor, 0)
	for _, s := range strs {
		actorSlice = appendIfNotExists(actorSlice, actor.Actor{
			Name:    s,
			Website: website,
		})
	}
	model.Create(actorSlice)
}

func splitter(r rune) bool {
	return r == ' ' || r == '-' || r == '_'
}
