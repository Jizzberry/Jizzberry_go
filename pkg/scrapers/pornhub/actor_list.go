package pornhub

import (
	"context"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"strconv"
	"strings"
	"time"
)

func (p Pornhub) ScrapeActorList(ctx context.Context) {
	c := colly.NewCollector(colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11"))

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*pornhub.*",
		Parallelism: 2,
		RandomDelay: 8 * time.Second,
	})

	if err != nil {
		helpers.LogError(err.Error(), p.GetWebsite())
	}

	q, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	c.OnHTML("body", func(element *colly.HTMLElement) {
		actorSlice := make([]actor.Actor, 0)
		element.ForEach("a[data-mxptype=Pornstar]", func(i int, element1 *colly.HTMLElement) {
			// Scrape only actor having full name
			split := strings.FieldsFunc(element1.Attr("data-mxptext"), splitter)
			if len(split) < 4 && len(split) > 1 {
				actorSlice = append(actorSlice, actor.Actor{
					Name:    element1.Attr("data-mxptext"),
					UrlID:   "",
					Website: p.GetWebsite(),
				})
			}
		})
		model := actor.Initialize()
		defer model.Close()

		model.Create(actorSlice)

	})

	c.OnError(func(response *colly.Response, err error) {
		helpers.LogError(err.Error(), p.GetWebsite())
	})

	c.OnRequest(func(request *colly.Request) {
		select {
		case <-ctx.Done():
			request.Abort()
			return
		default:
			break
		}
	})

	for i := 1; i < 1754; i++ {
		url := "https://www.pornhub.com/pornstars?o=t&page=" + strconv.Itoa(i)
		err := q.AddURL(url)
		if err != nil {
			helpers.LogError(err.Error(), p.GetWebsite())
		}
	}
	err = q.Run(c)
	if err != nil {
		helpers.LogError(err.Error(), p.GetWebsite())
	}
}

func splitter(r rune) bool {
	return r == ' ' || r == '.' || r == '-' || r == '_'
}
