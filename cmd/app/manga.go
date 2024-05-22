package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Manga struct {
	AniListID       int    `json:"aniListId"`
	ComicVineID     int    `json:"comicVineId"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	Publisher       string `json:"publisher"`
	Status          string `json:"status"`
	Year            int    `json:"year"`
	Description     string `json:"description"`
	NumberOfVolumes int    `json:"numberOfVolumes"`
	CoverImage      string `json:"coverImage"`
	URL             string `json:"url"`
}

type MangaDbRow struct {
	ID int `json:"Id"`
	Manga
}

type Volume struct {
	MangaID      int `json:"mangaId"`
	VolumeNumber int `json:"volumenumber"`
}

type VolumeDbRow struct {
	ID int `json:"id"`
	Volume
}

func addManga(c *gin.Context) {
	type Body struct {
		Title string `json:"title"`
	}
	var body Body

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := getAnilistData(body.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	publisher, err := parsePublisher(m.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, m)

	//// 1. Change Manga object, should be just title.
	//// 2. Parse Description for publisher
	// TODO: 3. ComicVine ?
	//// 4. Ids
	result, err := db.ExecContext(
		context.Background(),
		`INSERT INTO manga 
		(
			anilistid,
			comicvineid,
			title, 
			author, 
			publisher, 
			status, 
			year, 
			description, 
			numberofvolumes, 
			coverimage, 
			url
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.ID,
		0,
		m.Title.English,
		m.Staff.Nodes[0].Name.Full,
		publisher,
		m.Status,
		m.StartDate.Year,
		m.Description,
		m.Volumes,
		m.CoverImage.Large,
		m.SiteUrl,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, id)
}

func getManga(c *gin.Context) {
	var manga MangaDbRow
	id := c.Param("id")

	row := db.QueryRowContext(
		context.Background(),
		`SELECT * FROM manga WHERE id=?`, id,
	)
	err := row.Scan(
		&manga.ID,
		&manga.AniListID,
		&manga.ComicVineID,
		&manga.Title,
		&manga.Author,
		&manga.Publisher,
		&manga.Status,
		&manga.Year,
		&manga.Description,
		&manga.NumberOfVolumes,
		&manga.CoverImage,
		&manga.URL,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, manga)
}
