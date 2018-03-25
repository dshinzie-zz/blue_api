package main

import (
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
    db, _ := gorm.Open("sqlite3", "./gorm.db")
    defer db.Close()
}

