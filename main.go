package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
)

var jobTypeLinks []string

var jobDetailsLinks []string

var jobIDs []string

var mobileJobs []jobDetails

type jobDetails struct {
	Type  string
	Title string
	// Role      string
	Salary    string
	Location  string
	Seniority string
	Mode      string
}

// create collectors
// defaultCol is used to scrape job ids/types
var defaultCollector = colly.NewCollector(
	colly.AllowedDomains("www.getonbrd.com"),
)
var jobCollector = defaultCollector.Clone()

func main() {
	fmt.Println("Scraper is running...")

	// run collectors
	getJobIDs()
	getJobs()

	defaultCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	defaultCollector.Visit("https://www.getonbrd.com")

	// start scrapping by job id/type
	for i := range jobTypeLinks {
		jobCollector.Visit(jobTypeLinks[i])
	}
}

// default collector get job IDs
func getJobIDs() {
	defaultCollector.OnHTML("div[class=main-container]", func(e *colly.HTMLElement) {
		jobIDs = e.ChildAttrs("div.jobs", "id")

		for i := range jobIDs {
			link := "https://www.getonbrd.com/jobs/" + jobIDs[i]
			jobTypeLinks = append(jobTypeLinks, link)
		}
	})
}

func getJobs() {
	jobCollector.OnHTML("ul[class=sgb-results-list]", func(e *colly.HTMLElement) {
		jobDetailsLinks := e.ChildAttrs("a[href]", "href")

		jobCollector.OnHTML("div[id=right-col]", func(e *colly.HTMLElement) {

			location := e.ChildText("span[itemprop=address] > a[href] > span.location > span.tooltipster")
			salary := e.ChildText("span[itemprop=baseSalary] > span.tooltipster-basic > strong")
			// replace new lines
			re := regexp.MustCompile(`\r?\n`)

			// get job details
			jobDetails := jobDetails{
				Type:      e.ChildText("h2.size2 > a[href]"),
				Title:     e.ChildText("span[itemprop=title]"),
				Salary:    re.ReplaceAllString(salary, " "),
				Location:  re.ReplaceAllString(location, " "),
				Seniority: e.ChildText("span[itemprop=qualifications]"),
				Mode:      e.ChildText("span[itemprop=employmentType]"),
			}

			fmt.Printf("Job details %+v\n", jobDetails)

			mobileJobs = append(mobileJobs, jobDetails)
		})

		for i := range jobDetailsLinks {
			jobCollector.Visit(jobDetailsLinks[i])
		}
	})
}
