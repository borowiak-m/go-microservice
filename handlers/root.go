package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Root struct {
	log *log.Logger
}

func NewRoot(log *log.Logger) *Root {
	return &Root{log}
}

func (rt *Root) ServeHTTP(reqW http.ResponseWriter, req *http.Request) {
	rt.log.Println("Root request received")
	dat, err := io.ReadAll(req.Body)
	if err != nil {
		log.Panic("Error on root request parsing body")
		http.Error(reqW, "Oooops", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(reqW, "Hello %s", dat)
}
