package main

import (
	"fmt"
	"net/http"
	"os"
)

func getSites() []string {
	links := os.Args[1:]
	return links
}

func getStatus(l string, c chan string) {
	resp, err := http.Get(l)
	if err != nil {
		fmt.Printf("could not GET %s..\n", l)
		c <- l
		return
	}
	fmt.Printf("%s responded with status code %d\n", l, resp.StatusCode)
	c <- l
}

func main() {
	links := getSites()

	c := make(chan string)
	for _, link := range links {
		go getStatus(link, c)
	}
	for l := range c {
		go func(link string) {
			getStatus(link, c)
		}(l)
	}
}
