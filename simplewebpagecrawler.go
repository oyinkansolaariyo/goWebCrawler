package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var uri = flag.String("uri", "https://google.com", "Url to crawl")

func main() {
	flag.Parse()
	fmt.Println("This is the uri", uri)
	resp, err := http.Get(*uri)
	fmt.Println("The following error occcured while trying to connect to the specified url", err)
	if err != nil {
		fmt.Println("Shuting down app because an error  has occured")
		os.Exit(1)
	}

	defer resp.Body.Close() // run it when the function ends
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}
