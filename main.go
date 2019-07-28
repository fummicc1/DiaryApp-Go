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
	ID         string `json:"id" db:"id" gorm:"primary_key"`
	Title      string `json:"title" db:"title"`
	Content    string `json:"content" db:"content"`
	PosterName string `json:"poster_name" db:"poster_name"`
	CreatedAt  string `json:"created_at" db:"created_at"`
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
		diary.CreatedAt = time.Now().String()
		diary.ID = xid.New().String()
		db.Create(&diary)
		c.JSON(200, diary)
	})
	r.GET("diary/load", func(c *gin.Context) {
		var diaries []Diary
		db.Find(&diaries)
		fmt.Println(diaries)
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
		db.Delete(&diary)
		c.String(200, "Diary(id=?) Delete!", id)
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
