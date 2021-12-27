package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=golang&l=%EC%84%9C%EC%9A%B8"

type extractedJob struct {
	id       string
	location string
	title    string
	salary   string
	summary  string
}

func main() {
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}
	for i := 0; i < totalPages; i++ {
		job := <-c
		jobs = append(jobs, job...)
	}
	writeJobs(jobs)
	fmt.Println("Done!")
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}
	err = w.Write(headers)
	checkErr(err)

	for _, job := range jobs {
		jobSlice := []string{job.id, job.title, job.location, job.salary, job.summary}
		err = w.Write(jobSlice)
		checkErr(err)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPage(page int, cc chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL
	if page != 0 {
		pageURL += "&start=" + strconv.Itoa(page*10)
	}
	fmt.Println(pageURL)
	resp, err := http.Get(pageURL)
	checkErr(err)
	checkCode(resp.StatusCode)

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkErr(err)
	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	cc <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("id")
	title := cleanString(card.Find("h2>span").Text())
	location := cleanString(card.Find("div pre").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := cleanString(card.Find(".job-snippet").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
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
