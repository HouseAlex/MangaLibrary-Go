package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserDbRow struct {
	ID int `json:"id"`
	User
}

type UserToVolume struct {
	UserID   int `json:"userId"`
	VolumeID int `json:"volumeId"`
}

type UserToVolumeDbRow struct {
	ID int `json:"id"`
	UserToVolume
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
