package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
)

var jobIDs []string

var jobTypeLinks []string
var jobDetailsLinks []string

var uxJobs []jobDetails
var programmingJobs []jobDetails
var dsJobs []jobDetails
var mobileJobs []jobDetails
var supportJobs []jobDetails
var marketingJobs []jobDetails
var devopsQAJobs []jobDetails
var opsManagementJobs []jobDetails
var salesJobs []jobDetails
var advertisingJobs []jobDetails
var agileJobs []jobDetails
var hrJobs []jobDetails

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
		fmt.Println("Default collector Visiting", r.URL.String())
	})

	jobCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Job collector Visiting", r.URL.String())
	})

	defaultCollector.Visit("https://www.getonbrd.com")

	// start scrapping by job id/type
	for i := range jobTypeLinks {
		jobCollector.Visit(jobTypeLinks[i])
	}

	//fmt.Println("UX jobs: ", uxJobs)
	//fmt.Println("Programming  jobs: ", programmingJobs)
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

			//fmt.Printf("Job details %+v\n", jobDetails)

			switch jobDetails.Type {
			case "Design / UX":
				uxJobs = append(uxJobs, jobDetails)
			case "Programming":
				programmingJobs = append(programmingJobs, jobDetails)
			case "Data Science / Analytics":
				dsJobs = append(dsJobs, jobDetails)
			case "Mobile Development":
				mobileJobs = append(mobileJobs, jobDetails)
			case "Customer Support":
				supportJobs = append(supportJobs, jobDetails)
			case "Digital Marketing":
				marketingJobs = append(marketingJobs, jobDetails)
			case "SysAdmin / DevOps / QA":
				devopsQAJobs = append(devopsQAJobs, jobDetails)
			case "Operations / Management":
				opsManagementJobs = append(opsManagementJobs, jobDetails)
			case "Sales":
				salesJobs = append(salesJobs, jobDetails)
			case "Advertising & Media":
				advertisingJobs = append(advertisingJobs, jobDetails)
			case "Innovation & Agile":
				agileJobs = append(agileJobs, jobDetails)
			case "People & HR":
				hrJobs = append(hrJobs, jobDetails)
			}
		})

		for i := range jobDetailsLinks {
			jobCollector.Visit(jobDetailsLinks[i])
		}
	})
}
