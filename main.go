package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

func getAlbums(c *gin.Context) {

	var albums []Album

	rows, err := db.Query("SELECT * FROM album")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func addAlbums(c *gin.Context) {

	var alb Album
	if err := c.BindJSON(&alb); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, err)
		return
	}
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}

	var resultAlb Album = Album{
		ID:     id,
		Title:  alb.Title,
		Artist: alb.Artist,
		Price:  alb.Price,
	}

	c.IndentedJSON(http.StatusCreated, resultAlb)
}

func connectToDatabase() {
	cfg := mysql.Config{
		AllowNativePasswords: true,
		User:                 "root",
		Passwd:               "senha12345",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "go_api",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
}

func main() {
	connectToDatabase()

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", addAlbums)
	router.Run("localhost:3000")
}
