//package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var jobIDs, jobTypeLinks, jobDetailsLinks []string

var uxJobs, programmingJobs, dsJobs, mobileJobs, supportJobs, marketingJobs, devopsQAJobs, opsManagementJobs,
	salesJobs, advertisingJobs, agileJobs, hrJobs []jobDetails

var totalJobs, totalUxJobs, totalProgrammingJobs, totalDsJobs, totalMobileJobs, totalSupportJobs, totalMarketingJobs,
	totalDevopsQAJobs, totalOpsManagementJobs, totalSalesJobs, totalAdvertisingJobs, totalAgileJobs, totalHrJobs int64

var avgUxSalary, avgProgrammingSalary, avgMobileSalary, avgSupportSalary, avgMarketingSalary,
	vagDevopsQASalary, avgOpsMangmentSalary, avgSalesSalary, avgAdvertisingSalary, avgAgileSalary, avgHrSalary MeanMedian

var avgDsSalary []float64

type MeanMedian struct {
	numbers []float64
}

type jobDetails struct {
	Type  string
	Title string
	// Role      string
	Salary    float64
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
	/* for i := range jobTypeLinks {
	jobCollector.Visit(jobTypeLinks[i])
	} */
	jobCollector.Visit(jobTypeLinks[2])
	jobCollector.Visit(jobTypeLinks[3])

	//fmt.Println("UX jobs: ", uxJobs)
	//fmt.Println("Programming  jobs: ", programmingJobs)
	//fmt.Println("total UX josbs: ", totalUxJobs)

	fmt.Println("******************** SUMMARY ******************** ")
	fmt.Println("total MObile josbs: ", totalMobileJobs)
	fmt.Println("total DS josbs: ", totalDsJobs)
	fmt.Println("------------------------------------------------ ")
	fmt.Println("total jobs: ", totalJobs)
	fmt.Println("avg DS salary: ", avgDsSalary)
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

			fmt.Println("salary --> ", salary)

			//TODO: refactor
			// replace new lines
			regxNewline := regexp.MustCompile(`\r?\n`)
			regxNumbers := regexp.MustCompile("[0-9]+")
			salary = strings.Replace(salary, ",", "", -1)
			salaries := regxNumbers.FindAllString(salary, -1)
			var avgSalary, minSalary, maxSalary float64
			var errMin, errMax error

			fmt.Println("len ::> ", len(salaries))

			if len(salaries) == 1 {
				minSalary, errMin = strconv.ParseFloat(salaries[0], 64)
				if errMin != nil {
					fmt.Println("Error parsing min salary")
					return
				}

				avgSalary = minSalary
			}

			if len(salaries) > 1 {
				fmt.Println("salaries [0]: ", salaries[0])
				fmt.Println("salaries [1] ", salaries[1])

				minSalary, errMin = strconv.ParseFloat(salaries[0], 64)
				if errMin != nil {
					fmt.Println("Error parsing min salary")
					return
				}

				maxSalary, errMax = strconv.ParseFloat(salaries[1], 64)
				if errMax != nil {
					fmt.Println("Error parsing max salary")
					return
				}

				avgSalary = (minSalary + maxSalary) / 2
			}

			// get job details
			jobDetails := jobDetails{
				Type:      e.ChildText("h2.size2 > a[href]"),
				Title:     e.ChildText("span[itemprop=title]"),
				Salary:    avgSalary,
				Location:  regxNewline.ReplaceAllString(location, " "),
				Seniority: e.ChildText("span[itemprop=qualifications]"),
				Mode:      e.ChildText("span[itemprop=employmentType]"),
			}

			fmt.Printf("Job details %+v\n", jobDetails)

			// TODO: find better validation
			if jobDetails.Type != "" {
				totalJobs++
			}

			switch jobDetails.Type {
			case "Design / UX":
				uxJobs = append(uxJobs, jobDetails)
				totalUxJobs++
				// avgUxSalary += avgSalary
			case "Programming":
				programmingJobs = append(programmingJobs, jobDetails)
				totalProgrammingJobs++
			case "Data Science / Analytics":
				dsJobs = append(dsJobs, jobDetails)
				totalDsJobs++

				if avgSalary > 0 {
					avgDsSalary = append(avgDsSalary, avgSalary)
				}

			case "Mobile Development":
				mobileJobs = append(mobileJobs, jobDetails)
				totalMobileJobs++
			case "Customer Support":
				supportJobs = append(supportJobs, jobDetails)
				totalSupportJobs++
			case "Digital Marketing":
				marketingJobs = append(marketingJobs, jobDetails)
				totalMarketingJobs++
			case "SysAdmin / DevOps / QA":
				devopsQAJobs = append(devopsQAJobs, jobDetails)
				totalDevopsQAJobs++
			case "Operations / Management":
				opsManagementJobs = append(opsManagementJobs, jobDetails)
				totalOpsManagementJobs++
			case "Sales":
				salesJobs = append(salesJobs, jobDetails)
				totalSalesJobs++
			case "Advertising & Media":
				advertisingJobs = append(advertisingJobs, jobDetails)
				totalAdvertisingJobs++
			case "Innovation & Agile":
				agileJobs = append(agileJobs, jobDetails)
				totalAgileJobs++
			case "People & HR":
				hrJobs = append(hrJobs, jobDetails)
				totalHrJobs++
			}
		})

		for i := range jobDetailsLinks {
			jobCollector.Visit(jobDetailsLinks[i])
		}

		fmt.Println("avgDsSalary ---> ", avgDsSalary)
		fmt.Println("total ds Jobs ---> ", totalDsJobs)

		// avgUxSalary = avgUxSalary / totalUxJobs
		//avgDsSalary = avgDsSalary / totalDsJobs
		mmType := MeanMedian{avgDsSalary}

		fmt.Println("Median ::> ", mmType.calcMedian())
	})
}

func (mm *MeanMedian) calcMedian(n ...float64) float64 {
	sort.Float64s(mm.numbers) // sort the numbers

	mNumber := len(mm.numbers) / 2

	if mm.IsOdd() {
		return mm.numbers[mNumber]
	}

	return (mm.numbers[mNumber-1] + mm.numbers[mNumber]) / 2
}

// check if the total of numbers is odd or even
func (mm *MeanMedian) IsOdd() bool {
	if len(mm.numbers)%2 == 0 {
		return false
	}

	return true
}
