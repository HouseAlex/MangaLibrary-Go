package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var db *sql.DB

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type UserDbRow struct {
	ID int `json:"id"`
	User
}

type Manga struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	Publisher       string `json:"publisher"`
	Status          string `json:"status"`
	Year            int    `json:"year"`
	Description     string `json:"description"`
	NumberOfVolumes int    `json:"numberofvolumes"`
	CoverImage      string `json:"coverimage"`
	URL             string `json:"url"`
}

type MangaDbRows struct {
	ID int `json:"id"`
	Manga
}

type Volume struct {
	MangaID      int `json:"mangaid"`
	VolumeNumber int `json:"volumenumber"`
}

type VolumeDbRows struct {
	ID int `json:"id"`
	Volume
}

type UserToVolume struct {
	UserID   int `json:"userid"`
	VolumeID int `json:"volumeid"`
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
	// DB initialization
	err := initDatabase("data/MangaLibrary.db")
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

	router.Run("localhost:8080")
}

/*func getManga(c *gin.Context) {
	//id := c.Param("id")
	c.IndentedJSON(http.StatusOK, mangas)
}*/

func addManga(c *gin.Context) {
	var m Manga

	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.ExecContext(
		context.Background(),
		`INSERT INTO manga 
		(title, author, publisher, status, year, description, numberofvolumes, coverimage, url) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.Title, m.Author, m.Publisher, m.Status, m.Year, m.Description, m.NumberOfVolumes, m.CoverImage, m.URL,
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
			UserID INTEGER NOT NULL REFERENCES Users(UserID),
			VolumeID INTEGER NOT NULL REFERENCES Volumes(VolumeID)
		);`,
	)
	if err != nil {
		return err
	}
	return nil
}
