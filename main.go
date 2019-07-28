package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Diary struct {
	gorm.Model
	Title      string `json:"title"`
	Content    string `json:"content"`
	PosterName string `json:"posterName"`
}

var db *gorm.DB
var diaries []Diary

func main() {
	InitDB()
	router := NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8090"
	}
	router.Run(":" + port)
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/diary/create", func(c *gin.Context) {
		var diary Diary
		c.BindJSON(&diary)
		db.Create(&diary)
		c.JSON(200, diary)
	})
	r.GET("diary/load", func(c *gin.Context) {
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
		db.Delete(&diary)
		c.String(200, "Diary(id=?) Delete!", id)
	})
	return r
}

func InitDB() {
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	defer db.Close()
	CheckErr(err)
	db.AutoMigrate(&Diary{})
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
