package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

type User struct {
	ID        uint   `json:”id”`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func main() {
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=blue_api sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&User{})
	r := gin.Default()

	r.GET("/", GetUsers)
	r.POST("/people", CreatePerson)

	r.Run(":8080")
}

func GetUsers(c *gin.Context) {
	var people []User
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}
}

func CreatePerson(c *gin.Context) {
	var person User
	c.BindJSON(&person)
	fmt.Println(&person)
	db.Create(&person)
	fmt.Println(person)
	c.JSON(200, person)
}
