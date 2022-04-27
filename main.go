package main

import (
	"log"

	"k8smanager/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// proxyString := "http://172.19.0.3/file"
	// proxyURL, _ := url.Parse(proxyString)

	// rawURL := "http://172.19.0.3:32224"
	// url, _ := url.Parse(rawURL)

	// transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	// client := &http.Client{Transport: transport}

	// request, _ := http.NewRequest("GET", url.String(), nil)

	// res, _ := client.Do(request)
	//fmt.Println(res)
	router.RouterHandle()
}
