package res

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	//"reddico.co.uk/rankswarm-function/helpers/syscache"
	"strings"
	"time"
)

type SerpsResponse struct {
	Organic []struct {
		Rank        int    `json:"rank"`
		Link        string `json:"link"`
		DisplayLink string `json:"display_link"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Extensions  []struct {
			Text   string `json:"text"`
			Inline bool   `json:"inline"`
		} `json:"extensions,omitempty"`
	} `json:"organic"`
	Features []string
}

type ReturnOrganic struct {
	Rank        int    `json:"position"`
	Link        string `json:"url"`
	Description string `json:"text"`
}

type SerpsReturn struct {
	Organic []ReturnOrganic `json:"serps"`
	Query   struct {        // first found for check url
		Rank        int    `json:"position"`
		Link        string `json:"url"`
		Description string `json:"text"`
		Features    string `json:"features"` // features that domain ranks for
	} `json:"query"`
	Results  []ReturnOrganic `json:"results"`  // all positiions for check url
	Features []string        `json:"features"` // all of the features
}

func GetSerpsData(keyword string, checkUrl string, country string) (SerpsReturn, error, string) {

	response := SerpsResponse{}
	returnR := SerpsReturn{}

	var data []byte

	//cacheError := syscache.RdsGet("GetSerpsData" + keyword + country, &data)
	var cacheError error = nil

	if cacheError != nil {

		toCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// request using proxy
		proxyURL, _ := url.Parse("http://lum-customer-reddico-zone-residential_serp:ugi9ska3olge@zproxy.lum-superproxy.io:22225")

		httpClient := http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					dialler := &net.Dialer{
						Timeout:   30 * time.Second,
						KeepAlive: 30 * time.Second,
						Deadline:  time.Now().Add(30 * time.Second),
					}
					return dialler.DialContext(toCtx, network, addr)
				},
				MaxIdleConns:    50,
				IdleConnTimeout: 50 * time.Second,
			},
		}

		resp, err := httpClient.Get("http://www.google.com/search?q=" + url.QueryEscape(keyword) + "&gl=" + country + "&num=100&pws=0&lum_json=1")
		if err != nil {
			log.Println("proxy call error")
			log.Println(err)
			return returnR, err, ""
		}

		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("ioutil read error")
			fmt.Println(err)
			return returnR, err, ""
		}

		_, err = io.Copy(ioutil.Discard, resp.Body)
		if err != nil {
			log.Println("ioutil discard body error")
			fmt.Println(err)
			return returnR, err, ""
		}

		err = resp.Body.Close()
		if err != nil {
			log.Println("close body error")
			fmt.Println(err)
			return returnR, err, ""
		}

		//syscache.RdsSet("GetSerpsData"+keyword+country, data, 8 * time.Hour)

	}

	//Map json top level keys to find serp features
	c := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &c)
	if err != nil {
		log.Println("unmarshal 1 error")
		fmt.Println(string(data))
		return returnR, err, string(data)
	}

	var rankedFeatures string
	for s, _ := range c {
		if s != "general" && s != "organic" && s != "pagination" && s != "related" {
			if strings.Contains(string(c[s]), checkUrl) {
				rankedFeatures += s + ","
				fmt.Println("FOUND IN " + s)
			}
			response.Features = append(response.Features, s)
		}
	}
	rankedFeatures = strings.TrimSuffix(rankedFeatures, ",")

	returnR.Features = response.Features

	//Unmarshal rest of data to struct
	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Println("unmarshal 2 error")
		fmt.Println(string(data))
		return returnR, err, string(data)
	}

	firstFound := true
	for _, v := range response.Organic {

		linkParts := strings.Split(v.Link, "#")

		serp := ReturnOrganic{v.Rank, linkParts[0], v.Description}

		if strings.Contains(serp.Link, checkUrl) {
			returnR.Results = append(returnR.Results, serp)

			if firstFound {

				returnR.Query.Link = serp.Link
				returnR.Query.Description = serp.Description
				returnR.Query.Rank = serp.Rank
				returnR.Query.Features = rankedFeatures
				firstFound = false
			}
		}

		returnR.Organic = append(returnR.Organic, serp)

	}

	return returnR, nil, ""
}
