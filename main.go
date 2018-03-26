package main

import (
	"fmt"
    "time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/appleboy/gin-jwt"
)

var db *gorm.DB
var err error

type User struct {
	ID        uint   `json:"id"`
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

    authMiddleware := &jwt.GinJWTMiddleware{
        Realm:      "test zone",
        Key:        []byte("secret key"),
        Timeout:    time.Hour,
        MaxRefresh: time.Hour,
        Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
            if len(AuthUser(userId, password)) == 1  {
                return userId, true
            }

            return userId, false
        },
        Authorizator: func(userId string, c *gin.Context) bool {
            if userId == "admin" {
                return true
            }

            return false
        },
        Unauthorized: func(c *gin.Context, code int, message string) {
            c.JSON(code, gin.H{
                "code":    code,
                "message": message,
            })
        },
        // TokenLookup is a string in the form of "<source>:<name>" that is used
        // to extract token from the request.
        // Optional. Default value "header:Authorization".
        // Possible values:
        // - "header:<name>"
        // - "query:<name>"
        // - "cookie:<name>"
        TokenLookup: "header:Authorization",
        // TokenLookup: "query:token",
        // TokenLookup: "cookie:token",

        // TokenHeadName is a string in the header. Default value is "Bearer"
        TokenHeadName: "Bearer",

        // TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
        TimeFunc: time.Now,
    }

    auth := r.Group("/auth")
    auth.Use(authMiddleware.MiddlewareFunc())
    {
        auth.GET("/refresh_token", authMiddleware.RefreshHandler)
        auth.GET("/hello", helloHandler)
    }


    r.POST("/login", authMiddleware.LoginHandler)
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
    if err := c.ShouldBindJSON(&person); err != nil {
        c.AbortWithStatus(404)
        fmt.Println(err)
    } else {
        db.Create(&person)
        c.JSON(200, person)
    }
}

func helloHandler(c *gin.Context) {
    claims := jwt.ExtractClaims(c)
    c.JSON(200, gin.H{
        "userID": claims["id"],
        "text":   "Hello World.",
    })
}

func AuthUser(username, password string) []User {
    var users []User
    db.Where("email = ? AND password = ?", username, password).First(&users)
    fmt.Println(users)
    return users
}
