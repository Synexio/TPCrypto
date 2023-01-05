package main

import (
	"io/ioutil"
	"net/http"
)

func QueryKraken(path string) interface{} {
	// Requête de l'API de CoinGecko
	resp, err := http.Get(KrakenAPI + path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}
