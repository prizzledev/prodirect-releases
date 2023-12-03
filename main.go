package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls "github.com/bogdanfinn/tls-client"
	tlsProfiles "github.com/bogdanfinn/tls-client/profiles"
)

var (
	searchTerm = "dunk"      // eg. "dunk", "jordan", "new+balance"
	store      = "sportengb" // change to "soccerengb"
	searchURL  = fmt.Sprintf("https://query.published.live1.suggest.eu1.fredhopperservices.com/pro_direct/json?scope=//catalog01/en_GB/categories>{%s}&search=%s&callback=jsonpResponse", store, searchTerm)
)

func main() {
	_ = http.Client{}

	options := []tls.HttpClientOption{
		tls.WithClientProfile(tlsProfiles.Chrome_117),
	}

	client, err := tls.NewHttpClient(nil, options...)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`jsonpResponse\((\{.*?\})\)`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) == 0 {
		panic("no matches")
	}

	var pdsResp pdsResp
	err = json.Unmarshal([]byte(matches[1]), &pdsResp)
	if err != nil {
		panic(err)
	}

	for _, product := range pdsResp.SuggestionGroups[1].Suggestions {
		fmt.Println("----------")
		fmt.Println(product.Name)
		// price in pounds
		fmt.Println(fmt.Sprintf("£%.2f", product.CurrentPrice))
		fmt.Println(product.ProductURL)
		fmt.Println(product.ThumbURL)
		fmt.Println(convertFormatToDate(product.LaunchDate, product.LaunchTimeDelta))
	}
}

// funcs

func convertFormatToDate(launchDate string, delta int64) time.Time {
	year, _ := strconv.Atoi(launchDate[0:4])
	month, _ := strconv.Atoi(launchDate[4:6])
	day, _ := strconv.Atoi(launchDate[6:8])

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	deltaNs := delta*60*int64(time.Second) + (946080300+604800)*int64(time.Second)
	adjustedTime := time.Unix(0, deltaNs).UTC()

	date = date.Add(time.Hour*time.Duration(adjustedTime.Hour()) + time.Minute*time.Duration(adjustedTime.Minute()))
	return date
}

// types
type pdsResp struct {
	SuggestionGroups []SuggestionGroup `json:"suggestionGroups"`
}

type SuggestionGroup struct {
	IndexName   string       `json:"indexName"`
	IndexTitle  string       `json:"indexTitle"`
	Suggestions []Suggestion `json:"suggestions"`
}

type Suggestion struct {
	SearchTerm      string  `json:"searchterm,omitempty"`
	Name            string  `json:"name,omitempty"`
	SecondID        string  `json:"secondId,omitempty"`
	ThumbURL        string  `json:"_thumburl,omitempty"`
	CurrentPrice    float64 `json:"currentprice,string,omitempty"`
	PreviousPrice   float64 `json:"previousprice,string,omitempty"`
	ProductURL      string  `json:"producturl,omitempty"`
	SCProductURL    string  `json:"scproducturl,omitempty"`
	QuickRef        string  `json:"quickref,omitempty"`
	LaunchTimeDelta int64   `json:"launchtimedelta,string,omitempty"`
	LaunchDate      string  `json:"launchdate,omitempty"`
}
