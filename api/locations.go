package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	baseApiURL = "https://pokeapi.co/api"
	apiVersion = "v2"
	endPoint   = "location-area"
)

type LocationPayload struct {
	Id     int
	Name   string
	Limit  int
	Offset int
}

type Location struct {
	name string
	url  string
}

type Response struct {
	count    int
	next     string
	previous string
	results  []Location
}

func GetLocationsHandler(p *LocationPayload) {
	apiUrl := strings.Join([]string{baseApiURL, apiVersion, endPoint}, "/")

	params := url.Values{}
	params.Add("limit", strconv.Itoa(p.Limit))
	params.Add("offset", strconv.Itoa(p.Offset))

	res, err := http.Get(apiUrl + "?" + params.Encode())

	if err != nil {
		log.Printf("Request Failed: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	bodyString := string(body)
	log.Printf("Response: %s", bodyString)

	p.Offset += p.Limit
	// ResponseWithJSON(w, 200, res)
}
