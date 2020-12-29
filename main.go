package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

var jobsLink []string

// create collectors
// defaultCol is used to scrape job ids/types
var defaultCol = colly.NewCollector(
	colly.AllowedDomains("www.getonbrd.com"),
)

// is this really necessary? can't a colly collector be "reused"?
var uxCol = defaultCol.Clone()
var programmingCol = defaultCol.Clone()
var dsCol = defaultCol.Clone()
var mobileCol = defaultCol.Clone()

func main() {

	fmt.Println("Scraper is running...")

	// run collectors
	defaultCollector()
	uxCollector()
	programmingCollector()
	dsCollector()
	mobileCollector()

	defaultCol.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	uxCol.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	programmingCol.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	defaultCol.Visit("https://www.getonbrd.com")

	// start scrapping by job id/type
	uxCol.Visit(jobsLink[0])
	programmingCol.Visit(jobsLink[1])
	dsCol.Visit(jobsLink[2])
	mobileCol.Visit(jobsLink[3])

}

func defaultCollector() {
	defaultCol.OnHTML("div[class=main-container]", func(e *colly.HTMLElement) {
		jobIDs := e.ChildAttrs("div[id]", "id")

		for i := range jobIDs {
			fmt.Println("ids ---->", jobIDs[i])
			link := "https://www.getonbrd.com/jobs/" + jobIDs[i]
			jobsLink = append(jobsLink, link)
		}
	})
}

func uxCollector() {
	uxCol.OnHTML("ul[class=sgb-results-list]", func(e *colly.HTMLElement) {
		//uxLinks := e.ChildAttrs("a[href]", "href")
		//fmt.Println("uX links;:: ", uxLinks)
	})
}

func programmingCollector() {
	programmingCol.OnHTML("ul[class=sgb-results-list]", func(e *colly.HTMLElement) {
		//uxLinks := e.ChildAttrs("a[href]", "href")
		//fmt.Println("programming links;:: ", uxLinks)
	})
}

func dsCollector() {
	dsCol.OnHTML("ul[class=sgb-results-list]", func(e *colly.HTMLElement) {
		dsLinks := e.ChildAttrs("a[href]", "href")
		fmt.Println("data science links;:: ", dsLinks)
	})
}

func mobileCollector() {
	mobileCol.OnHTML("ul[class=sgb-results-list]", func(e *colly.HTMLElement) {
		mobileLinks := e.ChildAttrs("a[href]", "href")
		fmt.Println("mobile links;:: ", mobileLinks)
	})
}
