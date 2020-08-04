package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	url string
)

func main() {
	getURL()
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{url},
		ParseFunc: parse,
	}).Start()
}

func getURL() {
	fmt.Println("Enter the web-page link: ")
	_, err := fmt.Scanf("%s\n", &url)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`^http`)
	if !re.MatchString(url) {
		url = "http://" + url
	}
	re = regexp.MustCompile(`/$`)
	if !re.MatchString(url) {
		url += "/"
	}
	return
}

func parse(g *geziyor.Geziyor, r *client.Response)  {
	os.Mkdir("images", 0755)
	num1, num2 := 0, 0
	r.HTMLDoc.Find("img").Each(func(i int, s *goquery.Selection) {
		for _, v := range s.Nodes[0].Attr {
			if v.Key == "src" {
				err := downloadImage(v.Val, num1)
				if err != nil {
					num2++
				} else {
					num1++
				}
			}
		}
	})
	fmt.Printf("Images found: %d\n", num1 + num2)
	fmt.Printf("Images downloaded: %d\n", num1)
}

func downloadImage(path string, n int) error {
	re := regexp.MustCompile(`^http`)
	if !re.MatchString(path) {
		path = url + path
	}
	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create("images/" + strconv.FormatInt(int64(n + 1), 10) + path[strings.LastIndex(path, ".") : ])
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}