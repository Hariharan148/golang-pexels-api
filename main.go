package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
		PhotosApi = "https://api.pexels.com/v1/"
		VideosApi = "https://api.pexels.com/videos/"
)

type Client struct {
	Token string
	hc http.Client
	remainingTime int32
}




func NewClient(token string) *Client{
	c := http.Client{}

	return &Client{Token: token, hc: c}
}


type SearchResults struct {
	Page int32 `json:"page`
	PerPage int32 `json:"per_page`
	TotalResults int32 `json:"total_Results`
	NextPage string	`json:"next-page`
	Photos []Photo `json: "photos"`
}


type Photo struct{
	Id int32 `json: "id"`
	Width int32 `json:"width"`
	Height int32 `json:"height"`
	Url string `json: "url"`
	Photographer string `json:"photographer`
	PhotographerUrl string `json:"photographer_url"`
	Src PhotoSource `json:"src"`
}


type PhotoSource struct {
	Original string `json:"original"`
	Large string `json:"large"`
	Large2x string `json:large2x"`
	Medium string `json:"medium"`
	Small string `json:"small"`
	Potrait string  `json:"potrait"`
	Square string `json:"square"`
	Landscape string  `json:"landscape"`
	Tiny string `json:"tiny"`
}


type CuratedResults struct {
	PerPage int32 `json:"per_page"`
	Page int32 `json:"page"`
	NextPage int32 `json:next_page`
	Photos []Photo `json:"photos"`
}

func (c *Client) SearchPhotos(query string, Perpage int, page int) (*SearchResults, error){
	url := fmt.Sprintf(PhotosApi + "/search?query=%s&per_page=%d&page=%d", query, Perpage, page)

	resp, err := c.requestDoWithAuth("GET", url)


	defer resp.Body.Close()

	data , err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	var result SearchResults

	err = json.Unmarshal(data, &result)
	return &result, err
}


func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error){
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.Token)
	resp, err := c.hc.Do(req)
	if err != nil{
		return resp, err
	}

	times, err := strconv.Atoi(resp.Header.Get("X-Ratedlimit-Remaining"))
	if err != nil {
		return resp, nil
	} else{
		c.remainingTime = int32(times)
	}

	return resp, nil
}


func (c *Client)CuratedPhotos(perPage, page int)(*CuratedResults, error){
	url := fmt.Sprintf(PhotosApi+"curated?per_page=%d&page=%d", perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CuratedResults
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) getPhoto(id int)(*Photo, error) {
	url := fmt.Sprint(PhotosApi+"/photos/%d", id)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Photo

	err = json.Unmarshal(data, &result)

	return &result , err
}

func (c *Client) GetRandomPhotos() (*Photo, error){
	rand.Seed(time.Now().Unix())

	randNum := rand.Intn(1001)

	results, err := c.CuratedPhotos(1, randNum)
	if err == nil && len(results.Photos) == 1 {
		return &results.Photos[0], nil
	}

	return nil, err
	
}




func main(){

	os.Setenv("PEXELSTOKEN", "563492ad6f91700001000001e801300199cb45aa96bd9d0c3340fb09")

	TOKEN := os.Getenv("PEXELSTOKEN")

	c := NewClient(TOKEN)

	result, err := c.SearchPhotos("peacock", 30, 1)

	if err != nil {
		fmt.Errorf("Search Error: %v", err)
	}

	if result.Page == 0{
		fmt.Errorf("No results found")
	}

	fmt.Println(result)
}