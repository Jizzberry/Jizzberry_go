package scrapers

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/ghodss/yaml"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const component = "ScraperParser"

var scrapers = make([]scraper, 0)

type VideoDetails struct {
	Name    string   `json:"name"`
	Actors  []string `json:"actors"`
	Tags    []string `json:"tags"`
	Url     string   `json:"url"`
	Website string   `json:"website"`
}

type Videos struct {
	Name    string
	Url     string
	Website string
}

type scraper struct {
	path        string
	StudioList  bool
	ActorList   bool
	Actor       bool
	QueryVideos bool
	Video       bool
}

func RegisterScrapers() {
	err := filepath.Walk("./scrapers/", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			data := ParseYaml(path)
			if data != nil {
				scrapers = append(scrapers, scraper{
					path:        path,
					StudioList:  func() bool { _, ok := data[helpers.ScraperStudioList]; return ok }(),
					ActorList:   func() bool { _, ok := data[helpers.ScraperActorList]; return ok }(),
					Actor:       func() bool { _, ok := data[helpers.ScraperActor]; return ok }(),
					QueryVideos: func() bool { _, ok := data[helpers.ScraperVideos]; return ok }(),
					Video:       func() bool { _, ok := data[helpers.ScraperSingleVideo]; return ok }(),
				})
			}
		}
		return nil
	})

	if err != nil {
		helpers.LogError(err.Error(), component)
		return
	}
}

func MatchWebsite(website string) (bool, int) {
	for i := range scrapers {
		data := ParseYaml(scrapers[i].path)
		if val, ok := data[helpers.ScraperWebsite]; ok {
			if val == website {
				return true, i
			}
		}
	}
	return false, -1
}

func ParseYaml(path string) (yamlMap map[string]interface{}) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
	err = yaml.Unmarshal(yamlFile, &yamlMap)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	return yamlMap
}

func makeVideoStruct(name string, link string, website string) (videos Videos) {
	videos.Name = name
	videos.Url = link
	videos.Website = website
	return
}

func scrapeItem(regex []*regexp.Regexp, replacer string, subselector []interface{}, attr string, absolute bool, e *colly.HTMLElement, dest *string, condition func(string) bool) {
	for _, i := range subselector {
		var name []string
		if attr == "" {
			name = e.ChildTexts(i.(string))
		} else {
			name = e.ChildAttrs(i.(string), attr)
		}

		if i.(string) == "*" {
			name = []string{e.Text}
		}
		for _, n := range name {
			for _, r := range regex {
				if r.MatchString(n) {
					split := strings.Split(replacer, ";")
					var value string
					if len(split) > 1 {
						value = strings.TrimSpace(regexp.MustCompile(split[0]).ReplaceAllString(n, split[1]))
					} else {
						value = strings.TrimSpace(n)
					}

					if condition(value) {
						if absolute {
							*dest = e.Request.AbsoluteURL(value)
						} else {
							*dest = value
						}
					}
					if *dest != "" {
						return
					}
				}
			}
		}
	}
}

func scrapeList(selector interface{}, data map[string]interface{}, headers []string, destinations *[][]string, e *colly.HTMLElement, condition func(string, int) bool) {
	e.ForEach(selector.(string), func(i int, element *colly.HTMLElement) {
		tmp := make([]string, 0)
		for i := range headers {
			var str string
			getDataAndScrape(data, headers[i], element, &str, func(s string) bool {
				return condition(s, i)
			})
			if str == "" {
				return
			}
			tmp = append(tmp, str)
		}
		if len(tmp) == len(headers) {
			for i := range tmp {
				(*destinations)[i] = append((*destinations)[i], tmp[i])
			}
		}
	})
}

func compileRegex(regex interface{}) (r []*regexp.Regexp) {
	if regex == nil {
		reg, err := regexp.Compile(".*")
		if err != nil {
			helpers.LogError(err.Error(), component)
			return
		}
		r = append(r, reg)
		return
	} else {
		for _, re := range regex.([]interface{}) {
			reg, err := regexp.Compile(re.(string))
			if err != nil {
				helpers.LogError(err.Error(), component)
				reg, err := regexp.Compile(".*")
				if err != nil {
					helpers.LogError(err.Error(), component)
					continue
				}
				r = append(r, reg)
				continue
			}
			r = append(r, reg)
		}
	}
	return
}

func parseUrl(base string, query string) string {
	return strings.ReplaceAll(base, "%QUERY", strings.ReplaceAll(query, " ", "%20"))
}

func getData(data map[string]interface{}, header string) (subSelector []interface{}, r []*regexp.Regexp, replacer string, attr string, absolute bool) {
	attr = safeCastString(safeSelectFromMap(safeMapCast(data[header]), helpers.YamlForEachAttr))
	subSelector = safeCastSliceString(safeSelectFromMap(safeMapCast(data[header]), helpers.YamlSelector))
	r = compileRegex(safeSelectFromMap(safeMapCast(data[header]), helpers.YamlStringRegex))
	replacer = safeCastString(safeSelectFromMap(safeMapCast(data[header]), helpers.YamlStringReplace))
	absolute = safeCastBool(safeSelectFromMap(safeMapCast(data[header]), "absolute"))
	return
}

func getDataAndScrape(data map[string]interface{}, header string, e *colly.HTMLElement, dest *string, condition func(string) bool) {
	subSelector, r, replacer, attr, absolute := getData(data, header)
	scrapeItem(r, replacer, subSelector, attr, absolute, e, dest, condition)
}

func appendIfNotExists(slice []actor.Actor, actor2 actor.Actor) []actor.Actor {
	for _, a := range slice {
		if a.Name == actor2.Name {
			return slice
		}
	}
	return append(slice, actor2)
}

func safeMapCast(item interface{}) map[string]interface{} {
	if item != nil {
		if casted, ok := item.(map[string]interface{}); ok {
			return casted
		}
	}
	return nil
}
func safeSelectFromMap(item map[string]interface{}, key string) interface{} {
	if item != nil {
		if val, ok := item[key]; ok {
			return val
		}
	}
	return nil
}

func safeCastSliceString(item interface{}) []interface{} {
	if item != nil {
		if casted, ok := item.([]interface{}); ok {
			return casted
		}
	}
	return nil
}

func safeCastString(item interface{}) string {
	if casted, ok := item.(string); ok {
		return casted
	}
	return ""
}

func safeConvertInt(item interface{}) int {
	if str := safeCastString(item); str != "" {
		num, err := strconv.Atoi(str)
		if err != nil {
			helpers.LogError(fmt.Sprintf("Failed to convert %v to int", item), component)
			return -1
		}
		return num
	}
	helpers.LogError(fmt.Sprintf("Failed to convert %v to int", item), component)
	return -1
}

func safeCastBool(item interface{}) bool {
	if casted, ok := item.(bool); ok {
		return casted
	}
	return false
}

func getColly(onHtml func(e *colly.HTMLElement)) (c *colly.Collector) {
	c = colly.NewCollector()

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	if err != nil {
		helpers.LogError(err.Error(), component)
		return nil
	}

	c.OnError(func(response *colly.Response, e error) {
		helpers.LogError(e.Error(), component)
	})

	c.OnHTML("body", onHtml)

	return
}

func getQueue() (q *queue.Queue) {
	q, _ = queue.New(
		2,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)
	return
}
