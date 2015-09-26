package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type req struct {
	s  string
	ch chan string
}

func compute(ch chan req) {

	tc := time.Tick(10 * time.Second)

	re := regexp.MustCompile("^(\\d+)([rgbs])(.*)$")
	r := "0"
	g := "0"
	b := "0"
	t := "0"
	for {
		select {
		case <-tc:
			fmt.Printf("\n")

		case rq := <-ch:
			s := rq.s

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
				case "s":
					t = a[1]
				}
				s = a[3]
			}
			rp := fmt.Sprintf("%s,%s,%s,%s", t, r, g, b)
			fmt.Printf("%s.\n", rp)
			rq.ch <- rp
		}
	}
}

func main() {

	ch := make(chan req)

	go compute(ch)

	http.Handle("/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/bar/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			rc := make(chan string)
			ch <- req{s: string(body), ch: rc}
			s := <-rc
			w.Write([]byte(s))
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
