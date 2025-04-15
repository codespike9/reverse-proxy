package main

import (
	"log"
	"net/http"
	"reverseproxy/proxy"
)

func main() {

	if err := proxy.LoadConfig("config.yaml"); err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	proxy.InitLogger()
	http.HandleFunc("/", proxy.ReverseProxyHandler)

	log.Println("Server started on :3002")
	if err := http.ListenAndServe(":3002", nil); err != nil {
		log.Fatal(err)
	}
}
