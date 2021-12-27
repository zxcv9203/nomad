package scrapper

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

type extractedJob struct {
	id       string
	location string
	title    string
	salary   string
	summary  string
}

func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/jobs?q=" + term
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseURL)
	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, c)
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
	c := make(chan []string)
	checkErr(err)
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}
	err = w.Write(headers)
	checkErr(err)

	for _, job := range jobs {
		go func(job extractedJob) {
			c <- []string{job.id, job.title, job.location, job.salary, job.summary}
		}(job)
		checkErr(err)
	}
	for range jobs {
		job := <-c
		err = w.Write(job)
		checkErr(err)
	}
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPage(page int, url string, cc chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := url
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
	title := CleanString(card.Find("h2>span").Text())
	location := CleanString(card.Find("div pre").Text())
	salary := CleanString(card.Find(".salary-snippet").Text())
	summary := CleanString(card.Find(".job-snippet").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}
func getPages(url string) int {
	pages := 0
	resp, err := http.Get(url)
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
