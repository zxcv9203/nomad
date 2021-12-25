package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=golang&l=%EC%84%9C%EC%9A%B8"

func main() {
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}

func getPage(page int) {
	pageUrl := baseURL
	if page != 0 {
		pageUrl += "&start=" + strconv.Itoa(page*10)
	}
	fmt.Println(pageUrl)
}

func getPages() int {
	pages := 0
	resp, err := http.Get(baseURL)
	checkErr(err)
	checkCode(resp.StatusCode)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkErr(err)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(code int) {
	if code != 200 {
		log.Fatalln("Request failed with Status : ", code)
	}
}
