package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/aerogo/http/client"
)

type StaffNodes struct {
	ID   int `json:"id"`
	Name struct {
		Full string `json:"full"`
	} `json:"name"`
}

type AlMedia struct {
	ID    int `json:"id"`
	Title struct {
		English string `json:"english"`
	} `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Volumes     int    `json:"volumes"`
	Staff       struct {
		Nodes []StaffNodes `json:"nodes"`
	} `json:"staff"`
	CoverImage struct {
		Large string `json:"large"`
	} `json:"coverImage"`
	Status    string `json:"status"`
	SiteUrl   string `json:"siteUrl"`
	StartDate struct {
		Year int `json:"year"`
	} `json:"startDate"`
}

func getAnilistData(mangaName string) (*AlMedia, error) {
	type Variables struct {
		Search string `json:"search"`
		Type   string `json:"type"`
	}

	// !AUTHOR WILL NEED TO BE GOTTEN ANOTHER WAY - most likely comicvine

	// Query Body
	body := struct {
		Query     string    `json:"query"`
		Variables Variables `json:"variables"`
	}{
		Query: `
				query ($search: String, $type: MediaType) {
					Media (search: $search, type: $type) {
						id
						title {
							english
						}
						description
						type
						volumes
						staff(page:1, perPage:1) {
						nodes {
							id
							name {
								full
							}
						}
						}
						coverImage {
							large
						}
						status
						siteUrl
						startDate {
							year
						}
					}
				}
		`,
		Variables: Variables{
			Search: mangaName,
			Type:   "MANGA",
		},
	}

	// Query Response
	response := new(struct {
		Data struct {
			Media *AlMedia `json:"media"`
		} `json:"data"`
	})

	err := anilistQuery(body, response)

	if err != nil {
		return nil, err
	}

	return response.Data.Media, nil
}

func anilistQuery(body interface{}, target interface{}) error {
	var headers = client.Headers{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	response, err := client.Post("https://graphql.anilist.co").Headers(headers).BodyJSON(body).EndStruct(target)

	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("status: %d\n%s", response.StatusCode(), response.String())
	}

	return nil
}

func parsePublisher(description string) (string, error) {
	pattern := `\(Source: ([^)]+)\)`

	// Compile RegEx
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	matches := re.FindStringSubmatch(description)
	if matches == nil {
		return "No publisher found.", nil
	}

	// Index 1 because of RegEx works
	return matches[1], nil
}
