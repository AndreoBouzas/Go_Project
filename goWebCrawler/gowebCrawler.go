package Webcrawler

import (
	"fmt"
	"net/http"
	"os"
)

func Webcrawler() {
	baseURL := "http://youtube.com/jsfunc"
	responsehttp, err := http.Get(baseURL)
	checkErr(err)

	fmt.Println(responsehttp)

}
func checkErr(err error) {
	if err != nil {
		fmt.Println("erro", err)
		os.Exit(1)
	}
}
