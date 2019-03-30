package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

var seedUrl = flag.String("uri", "http://google.com", "Initial url to start crawling from")

func main() {
	flag.Parse()
	foundUrls := make(map[string]bool)
	//we need 2 channels , one to broadcast the newly found urls
	//the second one is to notify us when we are done with trasversing the current page

	urlChannel := make(chan string) //unbuffered channel
	finishedChannel := make(chan bool)

	go crawlPage(seedUrl, urlChannel, finishedChannel)

	//subscribe  to each channel.
	// a select statement is one that allows a go routine wait on communiction processes
	select {
	case url := <-urlChannel:
		fmt.Println("received")
		fmt.Println(url)
		foundUrls[url] = true

	}

	fmt.Println("Found", len(foundUrls), "unique urls")

	for url, _ := range foundUrls {
		fmt.Println("Url", url)
	}

	defer close(urlChannel)
	<-finishedChannel
}

func getHref(token html.Token) (ok bool, href string) {

	for _, attribute := range token.Attr {
		if attribute.Key == "href" {
			fmt.Println("found")
			fmt.Println(attribute.Val)
			href = attribute.Val
			ok = true
		}

	}
	return
}

func crawlPage(url *string, chUrl chan string, chFinished chan bool) {
	resp, err := http.Get(*url)

	if err != nil {
		fmt.Println("The following error occcured while trying to connect to the specified url", err)
		fmt.Println("Coudn't crawl the url", url)
		return
	}

	defer func() {
		chFinished <- true // notify a finished state at the end of the function
	}()

	defer resp.Body.Close() // run it when the function ends

	//convert html to tokens, to make it parsable
	tokeniizedHtml := html.NewTokenizer(resp.Body)

	for {
		token := tokeniizedHtml.Next()
		switch token {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			tag := tokeniizedHtml.Token()
			isAnchor := tag.Data == "a"
			//check if it is an a tag
			if !isAnchor {
				continue
			}

			ok, url := getHref(tag)
			fmt.Println("ok", ok, "url", url)

			if !ok {
				continue
			}

			fmt.Println("broadcast", url)
			if strings.Index(url, "http") == 0 {
				chUrl <- url
			}

		}

	}

}
