package pornhub

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor_details"
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
				actorDetails.Birthday = strings.TrimSpace(strings.ReplaceAll(element.Text, "Born:", ""))
			}
			if strings.ToLower(j) == "birthplace:" {
				actorDetails.Birthplace = strings.TrimSpace(strings.ReplaceAll(element.Text, "Birthplace:", ""))
			}
			if strings.ToLower(j) == "height:" {
				actorDetails.Height = strings.TrimSpace(strings.ReplaceAll(element.Text, "Height:", ""))
			}
			if strings.ToLower(j) == "weight:" {
				actorDetails.Weight = strings.TrimSpace(strings.ReplaceAll(element.Text, "Weight:", ""))
			}
		}
	})
}

func getDetails(url string) actor_details.ActorDetails {
	c := colly.NewCollector()

	actorDetails := actor_details.ActorDetails{}

	c.OnHTML("body", func(e *colly.HTMLElement) {
		v1(e, &actorDetails)

		if (actor_details.ActorDetails{}) == actorDetails {
			v2(e, &actorDetails)
		}
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Println(err)
	}
	c.Wait()

	return actorDetails
}

func (p Pornhub) ScrapeActor(name string) actor_details.ActorDetails {
	name = strings.ReplaceAll(name, " ", "-")
	return getDetails("https://www.pornhub.com/pornstar/" + name)
}

func (p Pornhub) GetWebsite() string {
	return "pornhub"
}
