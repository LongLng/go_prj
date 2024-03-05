package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []movie
}

type movie struct {
	Title string
	Year  string
}

func crawl(month int, day int) {

	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)
	infoCollector := c.Clone()

	c.OnHTML(".mode-detail", func(element *colly.HTMLElement) {
		profileUrl := element.ChildAttr("div.lister-item-image > a", "href")
		profileUrl = element.Request.AbsoluteURL(profileUrl)
		err := infoCollector.Visit(profileUrl)
		if err != nil {
			return
		}
		c.OnHTML("a.lister-page-next", func(element *colly.HTMLElement) {
			nextPage := element.Request.AbsoluteURL(element.Attr("href"))
			err := c.Visit(nextPage)
			if err != nil {
				return
			}
		})
		infoCollector.OnHTML("#content-2-wide", func(element *colly.HTMLElement) {
			tmpProfile := star{}
			tmpProfile.Name = element.ChildText("h1.header > span.itemprop")
			tmpProfile.Photo = element.ChildAttr("#name-poster", "src")
			tmpProfile.JobTitle = element.ChildText("#name-job-categories > a > span.itemprop")
			tmpProfile.BirthDate = element.ChildAttr("#name-born-info time", "datetime")

			tmpProfile.Bio = strings.TrimSpace(element.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))

			element.ForEach("div.knownfor-title", func(_ int, element *colly.HTMLElement) {
				tmpMovie := movie{}
				tmpMovie.Title = element.ChildText("div.knowfor-title-role > a.knownfor-ellipsis")
				tmpMovie.Year = element.ChildText("div.knownfor-year > span.knownfor-ellipsis")
				tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
			})
			js, err := json.MarshalIndent(tmpProfile, "", "   ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(js))
		})
		c.OnRequest(func(request *colly.Request) {
			fmt.Println("Visitting:", request.URL.String())
		})
		infoCollector.OnRequest(func(request *colly.Request) {
			fmt.Println("visiting profile URL", request.URL.String())

		})
	})
	startUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	err := c.Visit(startUrl)
	if err != nil {
		return
	}

}

func main() {
	month := flag.Int("month", 1, "Month to fetch birthdays for ")
	day := flag.Int("day", 1, "Day to fetch birthdays for")
	flag.Parse()
	crawl(*month, *day)
}
