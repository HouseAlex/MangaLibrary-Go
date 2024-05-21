package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/aerogo/http/client"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var db *sql.DB

type User struct {
	Username  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserDbRow struct {
	ID int `json:"id"`
	User
}

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

type UserToVolume struct {
	UserID   int `json:"userId"`
	VolumeID int `json:"volumeId"`
}

type UserToVolumeDbRow struct {
	ID int `json:"id"`
	UserToVolume
}

// Test Data
/*
var users = []User{
	{Username: "Wallu", FirstName: "Alex", LastName: "House"},
	{Username: "Lurk390", FirstName: "Mahmoud", LastName: "Elbasiouny"},
}

var mangas = []Manga{
	{Title: "Berserk", Author: "Kentarou Miura", Publisher: "Dark Horse", Status: "Active", Year: 2000, Description: "blah blah", NumberOfVolumes: 12, CoverImage: "blah blah", URL: "https://anilist.co/manga/30002/Berserk"},
	{Title: "Chainsaw Man", Author: "Tatsuki Fujimoto", Publisher: "Dark Horse", Status: "Active", Year: 2019, Description: "blah blah", NumberOfVolumes: 16, CoverImage: "blah blah", URL: "https://anilist.co/manga/105778/Chainsaw-Man/"},
	{Title: "Vinland Saga", Author: "Makoto Yukimura", Publisher: "Dark Horse", Status: "Active", Year: 2016, Description: "blah blah", NumberOfVolumes: 14, CoverImage: "blah blah", URL: "https://anilist.co/anime/101348/Vinland-Saga/"},
}*/

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// DB initialization
	err = initDatabase("data/MangaLibrary.db")
	if err != nil {
		log.Fatal("error initializaing DB connection: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("error initializing DB connection: ping error: ", err)
	}
	fmt.Println("database initialized..")

	// Gin Router initialization
	router := gin.Default()
	//router.GET("/get-manga/:id", getManga)
	router.POST("/add-manga", addManga)
	router.POST("/add-user", addUser)
	router.GET("get-user/:id", getUser)
	router.GET("get-manga/:id", getManga)

	router.Run("localhost:8080")
}

/*func getManga(c *gin.Context) {
	//id := c.Param("id")
	c.IndentedJSON(http.StatusOK, mangas)
}*/

func addManga(c *gin.Context) {
	var manga Manga

	if err := c.ShouldBindJSON(&manga); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := getAniListData(manga.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// TODO: 1. Change Manga object, should be just title.
	// TODO: 2. Parse Description for publisher
	// TODO: 3. ComicVine ?
	// TODO: 4. Ids

	c.JSON(http.StatusOK, m)

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
		"blah blah",
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

func addUser(c *gin.Context) {
	var u User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.ExecContext(
		context.Background(),
		`INSERT INTO users 
		(username, firstname, lastname) 
		VALUES (?, ?, ?)`,
		u.Username, u.FirstName, u.LastName,
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

func getUser(c *gin.Context) {
	var user UserDbRow
	id := c.Param("id")

	row := db.QueryRowContext(
		context.Background(),
		`SELECT * FROM users WHERE id=?`, id,
	)
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, user)
}

func initDatabase(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS Users (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Username TEXT NOT NULL UNIQUE,
			FirstName TEXT NOT NULL,
			LastName TEXT NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS Manga (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			AniListID INTEGER NULL,
			ComicVineId INTEGER NULL,
			Title TEXT,
			Author TEXT,
			Publisher TEXT,
			Status TEXT,
		--    VolumeType TEXT,
			Year INTEGER,
			Description TEXT,
			NumberOfVolumes INTEGER,
			CoverImage TEXT,
			URL TEXT
		);
		
		CREATE TABLE IF NOT EXISTS Volumes (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			MangaID INTEGER NOT NULL REFERENCES MangaInfo(MangaID),
			VolumeNumber INTEGER NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS UserToVolumes (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			UserID INTEGER NOT NULL REFERENCES Users(UserID),
			VolumeID INTEGER NOT NULL REFERENCES Volumes(VolumeID)
		);`,
	)
	if err != nil {
		return err
	}
	return nil
}

/*-----------------------------------------

		 ANILIST GRAPHQL QUERY CODE

-----------------------------------------*/

var headers = client.Headers{
	"Content-Type": "application/json",
	"Accept":       "application/json",
}

func Query(body interface{}, target interface{}) error {
	response, err := client.Post("https://graphql.anilist.co").Headers(headers).BodyJSON(body).EndStruct(target)

	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("status: %d\n%s", response.StatusCode(), response.String())
	}

	return nil
}

type StaffNodes struct {
	ID   int `json:"id"`
	Name struct {
		Full string `json:"full"`
	} `json:"name"`
}

type Media struct {
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

func getAniListData(mangaName string) (*Media, error) {
	type Variables struct {
		Search string `json:"search"`
		Type   string `json:"type"`
	}

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
						staff(sort:ROLE, page:1, perPage:1) {
						nodes {
							id
							name {
								full
							}
							primaryOccupations
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
			Media *Media `json:"media"`
		} `json:"data"`
	})

	err := Query(body, response)

	if err != nil {
		return nil, err
	}

	return response.Data.Media, nil
}

/*
GRAPH QL explorer
{
  Media(search:"Berserk", type:MANGA){
    id
    title {
      english
    }
    description
    type
    volumes
    staff(sort:ROLE, page:1, perPage:1) {
      nodes {
        id
        name {
          full
          native
        }
      }
    }
    coverImage {
      large
    }
    status
    siteUrl
  }
}

*/
