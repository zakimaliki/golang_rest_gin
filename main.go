package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//  install gin doc https://gin-gonic.com/docs/quickstart/
// auto reload doc https://github.com/codegangsta/gin

type Article struct {
	gorm.Model
	Title string
	Slug  string `gorm:"unique_index"`
	Desc  string `gorm:"type:text"`
}

var DB *gorm.DB

func main() {
	var err error
	DB, err = gorm.Open("mysql", "root:@/golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer DB.Close()

	// Migrate the schema
	DB.AutoMigrate(&Article{})
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		articles := v1.Group("/article")
		{
			articles.GET("/", getHome)
			articles.GET("/:slug", getArticle)
			articles.POST("/", postArticle)
			articles.PUT("/:id", updateArticle)
			articles.DELETE("/:id", deleteArticle)
		}
	}

	router.Run() // listen and serve on 0.0.0.0:8080
}

func getHome(c *gin.Context) {
	items := []Article{}
	DB.Find(&items)
	c.JSON(200, gin.H{
		"status": "Berhasil",
		"data":   items,
	})
}

func getArticle(c *gin.Context) {
	slug := c.Param("slug")

	var item Article

	if DB.First(&item, "slug = ?", slug).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  "Error",
			"message": "record not found",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   item,
	})
}

func postArticle(c *gin.Context) {

	item := Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Slug:  slug.Make(c.PostForm("title")),
	}

	DB.Create(&item)

	c.JSON(202, gin.H{
		"status": "Berhasil",
		"data":   item,
	})
}

func updateArticle(c *gin.Context) {
	id := c.Param("id")

	newArticle := Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Slug:  slug.Make(c.PostForm("title")),
	}

	var item Article

	DB.Model(&item).Where("id = ?", id).Updates(&newArticle)

	c.JSON(202, gin.H{
		"status": "Berhasil",
		"data":   item,
	})
}

func deleteArticle(c *gin.Context) {
	id := c.Param("id")
	var item Article
	DB.Delete(Article{}, "id = ?", id)

	c.JSON(200, gin.H{
		"status": "berhasil",
		"data":   item,
	})
}
