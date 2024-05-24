package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID        int    `json:"id" db:"ID"`
	Username  string `json:"userName" db:"Username"`
	FirstName string `json:"firstName" db:"FirstName"`
	LastName  string `json:"lastName" db:"LastName"`
}

type UserToVolume struct {
	ID       int `json:"id" db:"ID"`
	UserID   int `json:"userId" db:"UserID"`
	VolumeID int `json:"volumeId" db:"VolumeID"`
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
	var user User
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

func addUsersVolumes(c *gin.Context) {
	type Body struct {
		MangaId int   `json:"mangaId"`
		Volumes []int `json:"volumes"`
		UserId  int   `json:"userId"`
	}
	var body Body

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	placeholders := make([]string, len(body.Volumes))
	args := make([]interface{}, len(body.Volumes))
	for i, id := range body.Volumes {
		placeholders[i] = "?"
		args[i] = id
	}

	// Create query string
	query := fmt.Sprintf("SELECT * FROM Volumes WHERE mangaId = (%d) and id in (%s)", body.MangaId, strings.Join(placeholders, ", "))

	rows, err := db.QueryContext(
		context.Background(),
		query,
		args...,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var vol Volume
		err := rows.Scan(&vol.ID, &vol.MangaID, &vol.VolumeNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userVolDb, err := getUserVolume(vol.ID, body.UserId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if userVolDb == nil {
			result, err := db.ExecContext(
				context.Background(),
				`INSERT INTO usertovolumes
				(
					userid,
					volumeid
				)
				VALUES (?, ?)`,
				body.UserId,
				vol.ID,
			)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}

			if result == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "result was empty."})
			}
		}
	}
}

func getUserVolume(volumeID int, userID int) (*UserToVolume, error) {
	var userVolume UserToVolume

	row := db.QueryRowContext(
		context.Background(),
		`SELECT * FROM usertovolumes WHERE userid = ? AND volumeid = ?`,
		userID,
		volumeID,
	)
	err := row.Scan(
		&userVolume.ID,
		&userVolume.UserID,
		&userVolume.VolumeID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &userVolume, nil
}
