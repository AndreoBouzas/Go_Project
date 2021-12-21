package goWebcrawler

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/steelx/extractlinks"
)

var (
	config = &tls.Config{
		InsecureSkipVerify: true,
	}
	transport = &http.Transport{
		TLSClientConfig: config,
	}
	NetClient = &http.Client{
		Transport: transport,
	}
	queue = make(chan string)

	hasVisited = make(map[string]bool)
)

func Webcrawler() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		fmt.Println("Missing URL, e.g. go-webscrapper http://js.org/")
		os.Exit(1)
	}

	baseURL := arguments[0]
	go func() {
		queue <- baseURL
	}()

	for href := range queue {
		if !hasVisited[href] && isSameDomain(href, baseURL) {
			crawlUrl(href)
		}

	}

}

func crawlUrl(href string) {
	hasVisited[href] = true
	fmt.Printf("Crawling url -> %v \n ", href)
	responsehttp, err := NetClient.Get(href)
	checkErr(err)
	defer responsehttp.Body.Close()

	links, err := extractlinks.All(responsehttp.Body)
	checkErr(err)

	for _, link := range links {
		absoluteURL := toFixedURL(link.Href, href)
		go func() {
			queue <- absoluteURL
		}()

	}

	responsehttp.Body.Close()

}

func isSameDomain(href, baseUrl string) bool {
	uri, err := url.Parse(href)
	if err != nil {
		return false
	}
	parentUri, err := url.Parse(baseUrl)
	if err != nil {
		return false
	}
	if uri.Host != parentUri.Host {
		return false
	}
	return true
}

func toFixedURL(href, baseURL string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	toFixedUri := base.ResolveReference(uri)

	return toFixedUri.String()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("erro", err)
		os.Exit(1)
	}
}
