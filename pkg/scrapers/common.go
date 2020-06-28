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
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const component = "ScraperParser"

var scrapers = make([]scraper, 0)

type VideoDetails struct {
	Name    string
	Actors  []string
	Tags    []string
	Url     string
	Website string
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

func scrapeItem(regex []*regexp.Regexp, replacer string, subselector []interface{}, attr string, e *colly.HTMLElement, dest interface{}, condition func(string) bool) {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		helpers.LogError("Destination not a ptr", component)
		return
	}
	for _, i := range subselector {
		var name []string
		if attr == "" {
			name = e.ChildTexts(i.(string))
		} else {
			name = e.ChildAttrs(i.(string), attr)
		}
		for _, n := range name {
			for _, r := range regex {
				if r.MatchString(n) {
					split := strings.Split(replacer, ";")
					base := v.Elem()
					var value string
					if len(split) > 1 {
						value = strings.TrimSpace(regexp.MustCompile(split[0]).ReplaceAllString(n, split[1]))
					} else {
						value = strings.TrimSpace(n)
					}

					if val := reflect.ValueOf(value); base.Kind() != val.Kind() || condition(value) {
						base.Set(reflect.ValueOf(value))
					}

					if base.String() != "" {
						return
					}
				}
			}
		}
	}
}

func scrapeList(regex []*regexp.Regexp, replacer string, selector interface{}, subSelector []interface{}, attr string, e *colly.HTMLElement, dest interface{}, condition func(string) bool) {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		helpers.LogError("Destination not a ptr", component)
		return
	}

	base := v.Elem()
	if base.Kind() != reflect.Slice {
		helpers.LogError("Destination not a slice", component)
		return
	}

	e.ForEach(selector.(string), func(i int, element *colly.HTMLElement) {
		for _, s := range subSelector {
			var name []string
			if attr == "" {
				name = element.ChildTexts(s.(string))
			} else {
				name = element.ChildAttrs(s.(string), attr)
			}
			for _, n := range name {
				for _, r := range regex {
					if r.MatchString(n) {
						split := strings.Split(replacer, ";")
						base := v.Elem()
						var value string
						if len(split) > 1 {
							value = strings.TrimSpace(regexp.MustCompile(split[0]).ReplaceAllString(n, split[1]))
						} else {
							value = strings.TrimSpace(n)
						}

						if condition(value) {
							base.Set(reflect.Append(base, reflect.ValueOf(value)))
						}
					}
				}
			}
			if base.String() != "" {
				break
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

func getData(data map[string]interface{}, header string) (selector interface{}, subSelector []interface{}, r []*regexp.Regexp, replacer string, attr string) {
	attr = safeCastString(safeSelectFromMap(safeMapCast(data[header]), helpers.YamlForEachAttr))
	selector = safeCastString(safeSelectFromMap(safeMapCast(data[header]), helpers.YamlForEach))
	subSelector = safeCastSliceInterface(safeSelectFromMap(safeMapCast(data[header]), helpers.YamlSelector))
	r = compileRegex(safeSelectFromMap(safeMapCast(data[header]), helpers.YamlStringRegex))
	replacer = safeCastString(safeSelectFromMap(safeMapCast(data[header]), helpers.YamlStringReplace))
	return
}

func getDataAndScrape(data map[string]interface{}, header string, e *colly.HTMLElement, dest interface{}, multiple bool, condition func(string) bool) {
	selector, subSelector, r, replacer, attr := getData(data, header)

	if multiple {
		scrapeList(r, replacer, selector, subSelector, attr, e, dest, condition)
	} else {
		scrapeItem(r, replacer, subSelector, attr, e, dest, condition)
	}
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

func safeCastSliceInterface(item interface{}) []interface{} {
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
		helpers.LogError(err.Error(), component)
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
