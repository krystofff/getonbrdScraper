package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type attributes struct {
	Title       string   `json:"title"`
	Functions   string   `json:"functions"`
	Benefits    string   `json:"benefits"`
	Desirable   string   `json:"desirable"`
	Remote      bool     `json:"remote"`
	RemoteMod   string   `json:"remote_modality"`
	RemoteZone  string   `json:"remote_zone"`
	Country     string   `json:"country"`
	Category    string   `json:"category_name"`
	Perks       []string `json:"perks"`
	MinSalary   int      `json:"min_salary"`
	MaxSalary   int      `json:"max_salary"`
	Modality    string   `json:"modality"`
	Seniority   string   `json:"seniority"`
	PublishedAt int      `json:"published_at"`
}

type jobDetails struct {
	Data []struct {
		Id         string     `json:"id"`
		Type       string     `json:"type"`
		Attributes attributes `json:"attributes"`
	} `json:"data"`
}

func main() {

	/* url := "https://www.getonbrd.com/api/v0/categories?per_page=10&page=1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body)) */
	url := "https://www.getonbrd.com/api/v0/categories/programming/jobs?"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

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
	fmt.Println(jd.Data[0])

	// fmt.Println(jobs)

}
