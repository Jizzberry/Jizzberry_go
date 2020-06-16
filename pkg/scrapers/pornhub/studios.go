package pornhub

import (
	"context"
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/studios"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"strconv"
	"time"
)

func (p Pornhub) ScrapeStudiosList(ctx context.Context) {
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

	studioModel := studios.Initialize()

	c.OnHTML("body", func(element *colly.HTMLElement) {
		element.ForEach(".channelGridWrapper", func(i int, element *colly.HTMLElement) {
			studs := make([]studios.Studio, 0)
			for _, s := range element.ChildTexts(".usernameLink") {
				studs = append(studs, studios.Studio{Studio: s})
			}
			studioModel.Create(studs)
		})
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

	for i := 1; i < 90; i++ {
		url := "https://www.pornhub.com/channels?o=rk&page=" + strconv.Itoa(i)
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
