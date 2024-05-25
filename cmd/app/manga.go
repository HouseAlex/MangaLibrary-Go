package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Manga struct {
	ID              int    `json:"Id" db:"ID"`
	AniListID       int    `json:"aniListId" db:"AniListID"`
	ComicVineID     int    `json:"comicVineId" db:"ComicVineId"`
	Title           string `json:"title" db:"Title"`
	Author          string `json:"author" db:"Author"`
	Publisher       string `json:"publisher" db:"Publisher"`
	Status          string `json:"status" db:"Status"`
	Year            int    `json:"year" db:"Year"`
	Description     string `json:"description" db:"Description"`
	NumberOfVolumes int    `json:"numberOfVolumes" db:"NumberOfVolumes"`
	CoverImage      string `json:"coverImage" db:"CoverImage"`
	URL             string `json:"url" db:"URL"`
	IsActive        bool   `json:"isActive" db:"IsActive"`
	CreatedOn       string `json:"createdOn" db:"CreatedOn"`
}

type Volume struct {
	ID           int    `json:"id" db:"ID"`
	MangaID      int    `json:"mangaId" db:"MangaID"`
	VolumeNumber int    `json:"volumeNumber" db:"VolumeNumber"`
	IsActive     bool   `json:"isActive" db:"IsActive"`
	CreatedOn    string `json:"createdOn" db:"CreatedOn"`
}

func addManga(c *gin.Context) {
	type Body struct {
		Name        string `json:"name"`
		Volumes     int    `json:"volumes"`
		ComicVineId int    `json:"comicVineId"`
		Publisher   string `json:"publisher"`
	}
	var body Body

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mangaDb, err := getMangaFromComicVine(body.ComicVineId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var id int64

	if mangaDb == nil {
		m, err := getAnilistData(body.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		currentTime := time.Now().Format("2006-01-02")

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
				url,
				isactive,
				createdon
			) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			m.ID,
			body.ComicVineId,
			m.Title.English,
			m.Staff.Nodes[0].Name.Full,
			body.Publisher,
			m.Status,
			m.StartDate.Year,
			m.Description,
			body.Volumes,
			m.CoverImage.Large,
			m.SiteUrl,
			true,
			currentTime,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// Retrieve new manga identifier.
		id, err = result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// Assign volumes for publication
		err = addVolumes(id, body.Volumes)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		id = int64(mangaDb.ID)
	}

	c.JSON(http.StatusOK, id)
}

func addVolumes(mangaId int64, volumes int) error {
	currentTime := time.Now().Format("2006-01-02")
	for i := 1; i <= volumes; i++ {
		result, err := db.ExecContext(
			context.Background(),
			`INSERT INTO volumes
			(
				mangaid,
				volumeNumber,
				isactive,
				created
			)
			VALUES (?, ?, ?, ?)`,
			mangaId, i, true, currentTime,
		)
		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if rows != 1 && err != nil {
			return err
		}
	}

	return nil
}

func getManga(c *gin.Context) {
	var manga Manga
	id := c.Param("mangaId")

	err := db.GetContext(
		context.Background(),
		&manga,
		`SELECT * FROM manga WHERE id=?`, id,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, manga)
}

func getMangaFromComicVine(cvId int) (*Manga, error) {
	var manga Manga

	err := db.GetContext(
		context.Background(),
		&manga,
		`SELECT * FROM manga WHERE comicvineId=?`,
		cvId,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &manga, nil
}

func getMangaVolumes(c *gin.Context) {
	var volumes []Volume
	MangaID := c.Param("mangaId")

	err := db.SelectContext(
		context.Background(),
		&volumes,
		`SELECT * FROM volumes WHERE mangaId = ?`,
		MangaID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, volumes)
}
