package main

import (
	"fmt"
	"log"
	"net/http"

	gutilhttp "github.com/ryutah/gutil/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hogehoge")
	})

	if err := gutilhttp.StartServerwithGracefulShutdown(":8080", mux, func() {
		fmt.Println("finish ...")
	}); err != nil {
		log.Fatal(err)
	}
}
