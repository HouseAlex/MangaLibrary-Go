package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	// Load .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// DB initialization
	err = initDatabase("../../data/MangaLibrary.db")
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
