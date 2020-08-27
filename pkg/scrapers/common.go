package scrapers

import (
	"context"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/ghodss/yaml"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

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
				// Determines available functionalities of scraper
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
		helpers.LogError(err.Error())
		return
	}
}

// Returns index of scraper if query website string matches
// Returns false, -1 if scraper does not exist or can not be parsed
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
		helpers.LogError(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &yamlMap)
	if err != nil {
		helpers.LogError(err.Error())
	}

	return yamlMap
}

// Scrapes a single value into destination if conditions match
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
			// Match all available regex since a single web page can have variations
			for _, r := range regex {
				if r.MatchString(n) {
					// Replace certain keywords to avoid extra text in scraped data
					split := strings.Split(replacer, ";")
					var value string
					if len(split) > 1 {
						value = strings.TrimSpace(regexp.MustCompile(split[0]).ReplaceAllString(n, split[1]))
					} else {
						value = strings.TrimSpace(n)
					}

					if condition(value) {
						// If scraped data is a url and is relative, convert it to absolute
						if absolute {
							*dest = e.Request.AbsoluteURL(value)
						} else {
							*dest = value
						}
					}
					// If data is found, no need to continue to next iteration
					if *dest != "" {
						return
					}
				}
			}
		}
	}
}

// Scrapes array of strings in cases where a single web page contains multiple scrapable elements
func scrapeList(selector interface{}, data map[string]interface{}, headers []string, destinations *[][]string, e *colly.HTMLElement, condition func(string, int) bool) {
	e.ForEach(selector.(string), func(i int, element *colly.HTMLElement) {
		tmp := make([]string, 0)
		for i := range headers {
			var str string

			// Condition avoids scraping of certain text for specific headers
			getDataAndScrape(data, headers[i], element, &str, func(s string) bool {
				return condition(s, i)
			})

			if str == "" {
				// If any of headers is empty, leave destination as it is
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
		// If regex isn't provided, match everything
		reg, err := regexp.Compile(".*")
		if err != nil {
			helpers.LogError(err.Error())
			return
		}
		r = append(r, reg)
		return
	} else {
		for _, re := range regex.([]interface{}) {
			reg, err := regexp.Compile(re.(string))
			if err != nil {
				// Match everything if can not parse regex
				helpers.LogError(err.Error())
				reg, err := regexp.Compile(".*")
				if err != nil {
					helpers.LogError(err.Error())
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
	attr = helpers.SafeCastString(helpers.SafeSelectFromMap(helpers.SafeMapCast(data[header]), helpers.YamlForEachAttr))
	subSelector = helpers.SafeCastSlice(helpers.SafeSelectFromMap(helpers.SafeMapCast(data[header]), helpers.YamlSelector))
	r = compileRegex(helpers.SafeSelectFromMap(helpers.SafeMapCast(data[header]), helpers.YamlStringRegex))
	replacer = helpers.SafeCastString(helpers.SafeSelectFromMap(helpers.SafeMapCast(data[header]), helpers.YamlStringReplace))
	absolute = helpers.SafeCastBool(helpers.SafeSelectFromMap(helpers.SafeMapCast(data[header]), "absolute"))
	return
}

func getDataAndScrape(data map[string]interface{}, header string, e *colly.HTMLElement, dest *string, condition func(string) bool) {
	subSelector, r, replacer, attr, absolute := getData(data, header)
	scrapeItem(r, replacer, subSelector, attr, absolute, e, dest, condition)
}

// #### Cast types without panicking ####

// Returns instance of colly after setting default callbacks
func getColly(ctx context.Context, onHtml func(e *colly.HTMLElement)) (c *colly.Collector) {
	c = colly.NewCollector()

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	if err != nil {
		helpers.LogError(err.Error())
		return nil
	}

	c.OnRequest(func(request *colly.Request) {
		if ctx != nil {
			select {
			case <-ctx.Done():
				request.Abort()
				return
			default:
				return
			}
		}
	})

	c.OnError(func(response *colly.Response, e error) {
		helpers.LogError(e.Error())
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
