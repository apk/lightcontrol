package main

import (
	"fmt"
	"log"
	"net/http"
	"html"
	"regexp"
	"io/ioutil"
)

func compute(ch chan string) {

	re := regexp.MustCompile("^(\\d+)([rgb])(.*)$")
	r := "0"
	g := "0"
	b := "0"
	for {
		s := <- ch

		for {
			a := re.FindStringSubmatch(s)
			if a == nil {
				break
			}
			switch a[2] {
			case "r": 
				r = a[1]
			case "g": 
				g = a[1]
			case "b": 
				b = a[1]
			}
			s=a[3]
		}
		fmt.Printf(",%s,%s,%s.\n", r, g, b)
	}
}


func main () {

	ch := make(chan string)

	go compute(ch)

	http.Handle("/", http.FileServer(http.Dir(".")));

	http.HandleFunc("/bar/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			ch <- string(body)
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
