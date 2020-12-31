package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

var jobsLink []string

var mobileJobs []jobDetails

type jobDetails struct {
	Type  string
	Title string
	// Role      string
	Salary    int
	Location  string
	Seniority string
	Mode      string
}

// create collectors
// defaultCol is used to scrape job ids/types
var defaultCollector = colly.NewCollector(
	colly.AllowedDomains("www.getonbrd.com"),
)

// each job type/id have its own collector
var uxCollector = defaultCollector.Clone()
var programmingCollector = defaultCollector.Clone()
var dsCollector = defaultCollector.Clone()
var mobileCollector = defaultCollector.Clone()

func main() {

	fmt.Println("Scraper is running...")

	// run collectors
	getJobIDs()
	getUxJobs()
	getProgrammingJobs()
	getDsJobs()
	getMobileJobs()

	defaultCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	uxCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	mobileCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	defaultCollector.Visit("https://www.getonbrd.com")

	// start scrapping by job id/type
	uxCollector.Visit(jobsLink[0])
	programmingCollector.Visit(jobsLink[1])
	dsCollector.Visit(jobsLink[2])
	mobileCollector.Visit(jobsLink[3])
}

// default collector get job IDs
func getJobIDs() {
	defaultCollector.OnHTML("div[class=main-container]", func(e *colly.HTMLElement) {
		jobIDs := e.ChildAttrs("div[id]", "id")

		for i := range jobIDs {
			link := "https://www.getonbrd.com/jobs/" + jobIDs[i]
			jobsLink = append(jobsLink, link)
		}
	})
}

func getUxJobs() {
	uxCollector.OnHTML("ul[class=sgb-results-list]", func(e *colly.HTMLElement) {
		//uxLinks := e.ChildAttrs("a[href]", "href")
		//fmt.Println("uX links;:: ", uxLinks)
	})
}

func getProgrammingJobs() {
	programmingCollector.OnHTML("ul[class=sgb-results-list]", func(e *colly.HTMLElement) {
		//uxLinks := e.ChildAttrs("a[href]", "href")
		//fmt.Println("programming links;:: ", uxLinks)
	})
}

func getDsJobs() {
	dsCollector.OnHTML("ul[class=sgb-results-list]", func(e *colly.HTMLElement) {
		//dsLinks := e.ChildAttrs("a[href]", "href")
		//fmt.Println("data science links;:: ", dsLinks)
	})
}

func getMobileJobs() {
	mobileCollector.OnHTML("ul[class=sgb-results-list]", func(e *colly.HTMLElement) {
		mobileLinks := e.ChildAttrs("a[href]", "href")

		jobDetailsCol := mobileCollector.Clone()

		jobDetailsCol.OnHTML("div[id=right-col]", func(e *colly.HTMLElement) {
			// get job details
			jobDetails := jobDetails{
				Type:      "Mobile Developer", // harcoded helper in case there are page/site changes
				Title:     e.ChildText("span[itemprop=title]"),
				Salary:    0,
				Location:  e.ChildText("span[itemprop=address] > a[href] > span.location > span.tooltipster > span[itemprop=addressLocality]"),
				Seniority: e.ChildText("span[itemprop=qualifications]"),
				Mode:      e.ChildText("span[itemprop=employmentType]"),
			}

			fmt.Printf("Job details %+v\n", jobDetails)

			mobileJobs = append(mobileJobs, jobDetails)
		})

		fmt.Println("mobile job details: ", mobileJobs)

		for i := range mobileLinks {
			jobDetailsCol.Visit(mobileLinks[i])
		}
	})
}
