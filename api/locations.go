package api

import "net/http"

func GetLocationsHandler(w http.ResponseWriter, r *http.Request) {
	defer res.Body.Close()

	resp, error := http.Get("https://pokeapi.co/api/v2/location-area")
	if err != nil {
		// Do something with the error
	}
}
