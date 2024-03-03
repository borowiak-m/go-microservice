package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	// basic handler
	http.HandleFunc("/root", func(reqW http.ResponseWriter, req *http.Request) {
		log.Println("Root request received")
		dat, err := io.ReadAll(req.Body)
		if err != nil {
			log.Panic("Error on root request parsing body")
			http.Error(reqW, "Oooops", http.StatusBadRequest)
			return
		}
		log.Println(string(dat))
		_, err = fmt.Fprintf(reqW, "hello %s", dat)
		if err != nil {
			log.Panic("Error writing back to request")
			http.Error(reqW, "Oooops", http.StatusBadRequest)
			return
		}
	})
	// test handler
	http.HandleFunc("/test", func(reqW http.ResponseWriter, req *http.Request) {
		log.Println("Test request received")
		dat, err := io.ReadAll(req.Body)
		if err != nil {
			log.Panic("Error on root request parsing body")
			http.Error(reqW, "Oooops", http.StatusBadRequest)
			return
		}
		log.Println(string(dat))
		_, err = fmt.Fprintf(reqW, "hello %s", dat)
		if err != nil {
			log.Panic("Error writing back to request")
			http.Error(reqW, "Oooops", http.StatusBadRequest)
			return
		}
	})
	// create a web server
	http.ListenAndServe(":9090", nil)
}
