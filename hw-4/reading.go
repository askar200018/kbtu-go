package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"
)

const SPECIAL_DIVIDER = "$"

type Article struct {
	Author string
	Name   string
	Body   string
	Time   string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ConvertTextToArticleBySpecialDivider(data string, divider string) Article {
	splittedData := strings.Split(data, divider)
	article := Article{
		Author: splittedData[0],
		Name:   splittedData[1],
		Body:   splittedData[2],
		Time:   ConvertTimestampToDate(splittedData[3]),
	}
	return article
}

func ConvertTimestampToDate(timestamp string) string {
	tUnix, err := strconv.ParseInt(timestamp, 10, 64)
	check(err)
	timeT := time.Unix(tUnix, 0)
	return timeT.Format("2006 Jan 02")
}

func main() {
	data, err := os.ReadFile("test.txt")
	check(err)

	article := ConvertTextToArticleBySpecialDivider(string(data), SPECIAL_DIVIDER)

	articleJSON, err := json.Marshal(article)
	check(err)
	os.Stdout.Write(articleJSON)
}
