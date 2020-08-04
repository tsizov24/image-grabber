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
)

func main() {
	url := getURL()
	fmt.Println(url)
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{url},
		ParseFunc: parse,
	}).Start()
}

func getURL() (url string) {
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
	num1, num2 := 0, 0
	r.HTMLDoc.Find("img").Each(func(i int, s *goquery.Selection) {
		for _, v := range s.Nodes[0].Attr {
			if v.Key == "src" {
				fmt.Println(v.Val)
				err := downloadImage(v.Val)
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

func downloadImage(path string) error {
	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(path[len(path)-10:])
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}