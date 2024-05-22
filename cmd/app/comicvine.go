package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aerogo/http/client"
)

type CvMedia struct {
}

func getComicVineSearchData(mangaName string) (*CvMedia, error) {

	// Query Response
	response := new(struct {
		Data struct {
			Media *CvMedia `json:"media"`
		} `json:"data"`
	})
	var key = os.Getenv("CV_API_KEY")

	uri := "search/?api_key=" + key + "&format=json&query={manga_name}&resources=volume"
	err := comicVineQuery(uri, response)

	if err != nil {
		return nil, err
	}

	return response.Data.Media, nil
}

func comicVineQuery(uri string, target interface{}) error {
	var headers = client.Headers{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		//"(KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	}

	var baseAddress = "https://comicvine.gamespot.com/api/"

	response, err := client.Post(baseAddress + uri).Headers(headers).EndStruct(target)
	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("status: %d\n%s", response.StatusCode(), response.String())
	}

	return nil
}
