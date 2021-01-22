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
	name      string `json:"name"`
	dimension string `json:"dimension"`
}

type Categories struct {
	Data []struct {
		Id string `json:"id"`
	} `json:"data"`
}

var jobCategories []string

func main() {
	client := &http.Client{}

	getCategories("GET", categoriesUrl, client)

	getJobDetails(client, jobCategories)
}

func getCategories(method string, url string, client *http.Client) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var cat Categories

	err = json.Unmarshal(body, &cat)
	if err != nil {
		fmt.Println("error ", err)
		panic(err)
	}

	for i := range cat.Data {
		jobCategories = append(jobCategories, cat.Data[i].Id)
	}
}

// get job details by category
func getJobDetails(client *http.Client, categories []string) {

	for i := range categories {
		url := fmt.Sprintf("https://www.getonbrd.com/api/v0/categories/%s/jobs?", categories[i])

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
		}
		res, err := client.Do(req)
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		var jd jobDetails

		errr := json.Unmarshal(body, &jd)
		if errr != nil {
			fmt.Println("error ", errr)
			panic(err)
		}
		// fmt.Println(jd.Data[0])
		fmt.Printf("Job details %+v\n", jd.Data)
		fmt.Println("total jobs by category:: ", len(jd.Data))
	}

}
