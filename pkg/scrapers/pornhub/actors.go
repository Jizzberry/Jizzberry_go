package pornhub

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor_details"
	"github.com/gocolly/colly/v2"
	"strings"
)

func v1(e *colly.HTMLElement, actorDetails *actor_details.ActorDetails) {
	actorDetails.Name = e.ChildText("h1[itemprop=name]")
	e.ForEach(".infoPiece", func(i int, element *colly.HTMLElement) {
		for _, j := range element.ChildTexts("span") {
			if strings.ToLower(j) == "birthday:" {
				actorDetails.Birthday = element.ChildText(".smallInfo")
			}
			if strings.ToLower(j) == "birth place:" || strings.ToLower(j) == "birthplace:" {
				actorDetails.Birthplace = element.ChildText(".smallInfo")
			}
			if strings.ToLower(j) == "height:" {
				actorDetails.Height = element.ChildText(".smallInfo")
			}
			if strings.ToLower(j) == "weight:" {
				actorDetails.Weight = element.ChildText(".smallInfo")
			}
		}
	})
}

func v2(e *colly.HTMLElement, actorDetails *actor_details.ActorDetails) {
	actorDetails.Name = e.ChildText(".name")
	e.ForEach(".infoPiece", func(i int, element *colly.HTMLElement) {
		for _, j := range element.ChildTexts("span") {
			if strings.ToLower(j) == "born:" {
				actorDetails.Birthday = strings.TrimSpace(strings.ReplaceAll(element.Text, j, ""))
			}
			if strings.ToLower(j) == "birthplace:" {
				actorDetails.Birthplace = strings.TrimSpace(strings.ReplaceAll(element.Text, j, ""))
			}
			if strings.ToLower(j) == "height:" {
				actorDetails.Height = strings.TrimSpace(strings.ReplaceAll(element.Text, j, ""))
			}
			if strings.ToLower(j) == "weight:" {
				actorDetails.Weight = strings.TrimSpace(strings.ReplaceAll(element.Text, j, ""))
			}
		}
	})
}

func getDetails(url string) (actor_details.ActorDetails, error) {
	c := colly.NewCollector()

	actorDetails := actor_details.ActorDetails{}

	c.OnHTML("body", func(e *colly.HTMLElement) {
		v1(e, &actorDetails)

		if (actor_details.ActorDetails{}) == actorDetails {
			v2(e, &actorDetails)
		}
	})

	var err error
	c.OnError(func(response *colly.Response, e error) {
		err = e
	})

	if err != nil {
		return actorDetails, err
	}

	err = c.Visit(url)
	if err != nil {
		return actorDetails, err
	}
	c.Wait()

	return actorDetails, nil
}

func (p Pornhub) ScrapeActor(name string) (actor_details.ActorDetails, error) {
	name = strings.ReplaceAll(name, " ", "-")
	details, err := getDetails("https://www.pornhub.com/pornstar/" + name)
	if err != nil {
		helpers.LogError(err.Error(), p.GetWebsite())
	}
	return details, err
}

func (p Pornhub) GetWebsite() string {
	return "pornhub"
}
