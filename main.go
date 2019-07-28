package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/xid"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Diary struct {
	ID                  string `json:"id" db:"id" gorm:"primary_key"`
	Title               string `json:"title" db:"title"`
	Content             string `json:"content" db:"content"`
	PosterName          string `json:"poster_name" db:"poster_name"`
	DemandDeletionCount int    `json:"demend_deletion_count"`
	CreatedAt           string `json:"created_at" db:"created_at"`
}

var db *gorm.DB

func main() {
	InitDB()
	router := NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	router.Run(":" + port)
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.POST("/diary/create", func(c *gin.Context) {
		var diary Diary
		err := c.BindJSON(&diary)
		CheckErr(err)
		currentTime := time.Now()
		currentTime.In(time.Local)
		diary.CreatedAt = currentTime.Format("2006/01/02")
		diary.ID = xid.New().String()
		db.Create(&diary)
		c.JSON(200, diary)
	})
	r.GET("diary/load", func(c *gin.Context) {
		var diaries []Diary
		db.Find(&diaries)
		c.JSON(200, diaries)
	})
	r.GET("diary/load/:id", func(c *gin.Context) {
		id := c.Param("id")
		var diary Diary
		db.Find(&diary, "id=?", id)
		c.JSON(200, diary)
	})
	r.DELETE("diary/delete/:id", func(c *gin.Context) {
		id := c.Param("id")
		var diary Diary
		db.Find(&diary, "id=?", id)
		diary.DemandDeletionCount += 1
		if diary.DemandDeletionCount > 10 {
			db.Delete(&diary)
		} else {
			db.Save(&diary)
		}
		c.String(200, " Delete Request for This Diary Sent.", id)
	})
	return r
}

func InitDB() {
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	db.AutoMigrate(&Diary{})
	CheckErr(err)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
