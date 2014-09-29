package main

import (
	"fmt"
	"net/http"
	"encoding/xml"
	"log"
	"os"
	"strings"
)

func checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Enclosures []struct {
		Url string `xml:"url,attr"`
	} `xml:"channel>item>enclosure"`
}

func gtfu(url string, resultchan chan<- string) {
	result := ""
	defer func () { resultchan<-result }()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	dec := xml.NewDecoder(resp.Body)

	var rss RSS
	if dec.Decode(&rss) != nil{
		fmt.Println(err)
		return
	}

	if len(rss.Enclosures) == 0 {
		return
	}

	result = rss.Enclosures[0].Url
}

func main() {
	resultchan := make(chan string, len(os.Args) - 1)
	for _, url := range os.Args[1:] {
		go gtfu(url, resultchan)
	}

	urls := make([]string, 0)
	for i := 1; i < len(os.Args); i++ {
		urls = append(urls, <-resultchan)
	}

	fmt.Println(strings.Join(urls, " "))
}
