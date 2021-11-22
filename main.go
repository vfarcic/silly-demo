package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", handler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Println("Starting the server on " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	output := "This is a silly demo"
	if len(req.URL.Query()["fail"]) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		output = "Something terrible happened"
	}
	version := os.Getenv("VERSION")
	if len(version) > 0 {
		output = fmt.Sprintf("%s version %s", output, version)
	}
	if len(req.URL.Query()["html"]) > 0 {
		output = fmt.Sprintf("<h1>%s</h1>", output)
	}
	output = fmt.Sprintf("%s\n", output)
	fmt.Fprintf(w, output)
}
