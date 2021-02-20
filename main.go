package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	categoriesUrl string = "https://www.getonbrd.com/api/v0/categories?per_page=10&page=1"
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

type jobCategory struct {
	categoryName string
	total        int
}

type avgSalByCategory struct {
	categoryName string
	total        int
}

type jobsByCategory []jobCategory

type avgSalariesByCat []avgSalByCategory

type jobsDetails []jobDetails

var (
	totalJobsByCategory   jobsByCategory
	jobCategories         []string
	totalJobs             int
	avgSalariesByCategory avgSalariesByCat
)

func main() {

	getJobCategories()

	getJobDetails()
	fmt.Printf("total jobs by category %+v\n", totalJobsByCategory)
	fmt.Println("Total jobs: ", totalJobs)

	fmt.Println("AVG salaries: ", avgSalariesByCategory)
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

		// fmt.Println("DATA:: ", jd.Data[0])

		avgSalariesByCategory = append(avgSalariesByCategory, avgSalByCategory{jobCategories[i], getSalaryByCategory(jd)})

		totalJobs += len(jd.Data)
		totalJobsByCategory = append(totalJobsByCategory, jobCategory{jobCategories[i], len(jd.Data)})
	}

}

func getSalaryByCategory(jobDetails jobDetails) int {

	var avgSalaryByCategory, jobsLen int

	for _, data := range jobDetails.Data {
		if data.Attributes.MinSalary > 0 || data.Attributes.MaxSalary > 0 {
			jobsLen++
			avgSalaryByCategory += (data.Attributes.MinSalary + data.Attributes.MaxSalary) / 2
		}
	}
	return avgSalaryByCategory / jobsLen

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
