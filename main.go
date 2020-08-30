package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func getStatus(l string, d chan string, c *gin.Context) {
	resp, err := http.Get(l)
	if err != nil {
		c.String(http.StatusOK, "could not GET %s\n", l)
		// fmt.Printf("could not GET %s\n", l)
		d <- l
		return
	}
	if resp.StatusCode == 200 {
		c.String(http.StatusOK, "%s is UP (%d)\n", l, resp.StatusCode)
		// fmt.Printf("%s is UP (%d)\n", l, resp.StatusCode)
		d <- l
	} else {
		c.String(http.StatusOK, "%s responded with (%d)\n", l, resp.StatusCode)
		// fmt.Printf("%s responded with (%d)\n", l, resp.StatusCode)
		d <- l
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/check", func(c *gin.Context) {
		link := c.Request.URL.Query().Get("url")
		links := []string{link}

		d := make(chan string)
		for _, link := range links {
			go getStatus(link, d, c)
		}
		for l := range d {
			go func(link string) {
				getStatus(link, d, c)
			}(l)
		}
	})
	router.Run(":" + port)
}
