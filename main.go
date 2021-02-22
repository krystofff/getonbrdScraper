package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

const (
	categoriesUrl string = "https://www.getonbrd.com/api/v0/categories?"
	tagsUrl       string = "https://www.getonbrd.com/api/v0/tags?page=2"
)

type attributes struct {
	Title     string `json:"title"`
	Functions string `json:"functions"`
	// Benefits    string   `json:"benefits"`
	Desirable  string `json:"desirable"`
	Remote     bool   `json:"remote"`
	RemoteMod  string `json:"remote_modality"`
	RemoteZone string `json:"remote_zone"`
	Country    string `json:"country"`
	Category   string `json:"category_name"`
	// Perks       []string `json:"perks"`
	MinSalary int    `json:"min_salary"`
	MaxSalary int    `json:"max_salary"`
	Modality  string `json:"modality"`
	Seniority string `json:"seniority"`
	// PublishedAt int      `json:"published_at"`
}

type jobDetails struct {
	Data []struct {
		Id string `json:"id"`
		// Type string `json:"type"`
		Attributes attributes `json:"attributes"`
	} `json:"data"`
}

type catAttributes struct {
	Name      string `json:"name"`
	Dimension string `json:"dimension"`
}

type Categories struct {
	Data []struct {
		Id string `json:"id"`
	} `json:"data"`
}

type Tags struct {
	Data []struct {
		Id string `json:"id"`
	} `json:"data"`
}

type jobTotalCategory struct {
	categoryName string
	total        int
}

var (
	jobCategories            []string
	totalJobs                int
	totalJobsByCategory      []jobTotalCategory
	avgSalariesByCategory    []jobTotalCategory
	medianSalariesByCategory []jobTotalCategory
	avgSalariesList          []int
)

func main() {

	/* 	getJobCategories()

	   	getJobDetails()
	   	fmt.Println("Total jobs: ", totalJobs)
	   	fmt.Println("--------------------------------------------------------")
	   	fmt.Printf("Total jobs by category: %+v\n", totalJobsByCategory)
	   	fmt.Println("--------------------------------------------------------")
	   	fmt.Printf("AVG salaries by category: %+v\n", avgSalariesByCategory)
	   	fmt.Println("--------------------------------------------------------")
	   	fmt.Printf("Median salaries by category: %+v\n", medianSalariesByCategory)
	   	fmt.Println("--------------------------------------------------------")
	*/getTags()
}

func getJobCategories() {
	body := getRequest(categoriesUrl)

	var cat Categories

	err := json.Unmarshal(body, &cat)
	if err != nil {
		fmt.Println("error ", err)
		panic(err)
	}

	for i := range cat.Data {
		jobCategories = append(jobCategories, cat.Data[i].Id)
	}
}

// get job details by category
func getJobDetails() {
	for i := range jobCategories {
		url := fmt.Sprintf("https://www.getonbrd.com/api/v0/categories/%s/jobs?", jobCategories[i])

		body := getRequest(url)

		var jd jobDetails

		err := json.Unmarshal(body, &jd)
		if err != nil {
			fmt.Println("error ", err)
			panic(err)
		}

		totalJobs += len(jd.Data)
		totalJobsByCategory = append(totalJobsByCategory, jobTotalCategory{jobCategories[i], len(jd.Data)})
		getSalaryData(jd)
	}
}

// gets average and median salary by category
func getSalaryData(jobDetails jobDetails) {

	var avgSalaryByCategory, jobsLen, totalAvgSalaryByCategory int
	avgSalariesList = nil

	for _, data := range jobDetails.Data {
		if data.Attributes.MinSalary > 0 || data.Attributes.MaxSalary > 0 {
			jobsLen++
			avgSalaryByCategory = (data.Attributes.MinSalary + data.Attributes.MaxSalary) / 2
			totalAvgSalaryByCategory += (data.Attributes.MinSalary + data.Attributes.MaxSalary) / 2
			avgSalariesList = append(avgSalariesList, avgSalaryByCategory)
		}
	}

	if jobsLen > 0 {
		avgSalariesByCategory = append(avgSalariesByCategory,
			jobTotalCategory{jobDetails.Data[0].Attributes.Category, totalAvgSalaryByCategory / jobsLen})
	}

	if len(avgSalariesList) > 0 {
		medianSalariesByCategory = append(medianSalariesByCategory,
			jobTotalCategory{jobDetails.Data[0].Attributes.Category, getMedian(avgSalariesList)})
	}
}

func getMedian(n []int) int {
	sort.Ints(n)          // sort the numbers
	mNumber := len(n) / 2 // middle number, truncates if odd

	// is Odd?
	if len(n)%2 != 0 {
		return n[mNumber]
	}

	return (n[mNumber-1] + n[mNumber]) / 2
}

// TODO: should get tags by page: max = 100 per page
func getTags() {
	body := getRequest(tagsUrl)

	var tag Tags

	err := json.Unmarshal(body, &tag)
	if err != nil {
		fmt.Println("error ", err)
		panic(err)
	}

	fmt.Println("LEN --> ", len(tag.Data))
	for _, data := range tag.Data {
		fmt.Println("Tag:: ", data)
	}
}

func getRequest(url string) []byte {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
	}

	return body
}
