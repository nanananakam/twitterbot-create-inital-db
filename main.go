package main

import (
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
	"os"
)

type Tweet struct {
	gorm.Model
	TwitterID string
	Tweet     string
}

func main() {

	db, err := gorm.Open("sqlite3", "tweets.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&Tweet{})
	db.Model(&Tweet{}).AddUniqueIndex("idx_twitter_id", "Tweet")

	fp, err := os.Open("tweets.csv")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true

	tx := db.Begin()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		tweet := Tweet{
			TwitterID: record[0],
			Tweet: record[1],
		}
		fmt.Println(tweet)
		tx.Create(&tweet)
	}

	tx.Commit()

	var tweets []Tweet

	db.Find(&tweets)

	for _, tweet := range tweets {
		fmt.Println(tweet)
	}
}
