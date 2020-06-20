package pornhub

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/gocolly/colly/v2"
	"path/filepath"
	"strconv"
	"strings"
)

func getImages2(url string, actorId int64) error {
	c := colly.NewCollector()

	gotImg := false

	c.OnHTML("body", func(e *colly.HTMLElement) {
		err := downloadImage(e.ChildAttr("img[id=getAvatar]", "src"), actorId)
		if err == nil {
			gotImg = true
		}
	})

	err := c.Visit(url)
	if err != nil {
		return err
	}
	c.Wait()

	if !gotImg {
		return fmt.Errorf("failed to aquire image")
	}

	return nil
}

func getImages(url string, actorId int64) error {
	c := colly.NewCollector()

	gotImg := false

	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach(".pornstarAlbumListBlock", func(i int, element *colly.HTMLElement) {
			if i == 1 {
				err := downloadImage(element.ChildAttr("img", "src"), actorId)
				if err == nil {
					gotImg = true
				}
			}
		})
	})

	err := c.Visit(url)
	if err != nil {
		return err
	}
	c.Wait()

	if !gotImg {
		return fmt.Errorf("failed to aquire image for actorID %d", actorId)
	}

	return nil
}

func (p Pornhub) ScrapeImage(name string, actorId int64) {
	name = strings.ReplaceAll(name, " ", "-")
	err := getImages("https://www.pornhub.com/pornstar/"+name+"/official_photos", actorId)
	if err != nil {
		err = getImages2("https://www.pornhub.com/pornstar/"+name, actorId)
		if err != nil {
			helpers.LogError(err.Error(), p.GetWebsite())
		}
	}
}

func downloadImage(url string, actorId int64) error {
	imageC := colly.NewCollector()
	var err error
	imageC.OnResponse(func(response *colly.Response) {
		err = response.Save(filepath.FromSlash(helpers.ThumbnailPath + "/p" + strconv.FormatInt(actorId, 10)))
	})
	err = imageC.Visit(url)
	return err
}
