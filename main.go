package main

import (
	"encoding/csv"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/shogo82148/go-mecab"
	"io"
	"os"
	"regexp"
)

type Tweet struct {
	gorm.Model
	TwitterID string `gorm:"unique_index"`
	Tweet     string `gorm:"type:varchar(1024)"`
}

type Words struct {
	gorm.Model
	Word1 string `gorm:"index;type:varchar(1024)"`
	Word2 string `gorm:"type:varchar(1024)"`
}

func main() {

	db, err := gorm.Open("sqlite3", "tweets.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db2, err := gorm.Open("sqlite3", "words.db")
	if err != nil {
		panic(err)
	}
	defer db2.Close()

	db.AutoMigrate(&Tweet{})

	db2.AutoMigrate(&Words{})

	fp, err := os.Open("tweets.csv")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = ','
	reader.LazyQuotes = true

	tagger, err := mecab.New(map[string]string{"output-format-type": "wakati"})
	if err != nil {
		panic(err)
	}
	defer tagger.Destroy()

	tx := db.Begin()
	tx2 := db2.Begin()

	filterRep := regexp.MustCompile(`(RT|@[^ 　]+|http[^ 　]+|\\)`)
	ngRep := regexp.MustCompile(`(死ね|殺|爆破)`)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		tweet := Tweet{
			TwitterID: record[0],
			Tweet:     record[1],
		}
		tx.Create(&tweet)
		if !ngRep.MatchString(record[1]) {
			tweetString := filterRep.ReplaceAllString(record[1], "")
			fmt.Println(tweetString)
			node, err := tagger.ParseToNode(tweetString)
			if err != nil {
				panic(err)
			}
			var word1 string
			var word2 string
			for ; node != (mecab.Node{}); node = node.Next() {
				if node.Surface() != "" {
					word2 = node.Surface()
					words := Words{
						Word1: word1,
						Word2: word2,
					}
					tx2.Create(&words)
					word1 = word2
				}
			}
			if word1 != "" {
				word2 = ""
				words := Words{
					Word1: word1,
					Word2: word2,
				}
				tx2.Create(&words)
			}
		}

	}

	tx.Commit()
	tx2.Commit()
}
