package main

import (
	"encoding/json"
	"fmt"
	"github.com/bcampbell/fuzzytime"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Times struct {
	Unix    interface{} `json:"unix"`
	Natural interface{} `json:"natural"`
}

var natural string
var unix int64

func Router(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" || r.URL.Path[1:] == "favicon.ico" {
		Start(w, r)
	} else if r.URL.Path[1:] != "" {
		var times = Times{}
		i, err := strconv.ParseInt(r.URL.Path[1:], 10, 64)
		if err == nil { // URL is an integer
			unix = i
			date := time.Unix(unix, 0)
			natural := date.Format("January 2, 2006")
			times = Times{Unix: unix, Natural: natural}
		} else { // URL is not an integer
			urlAfter, _ := url.QueryUnescape(r.URL.Path[1:])
			str, _, _ := fuzzytime.Extract(urlAfter)
			if str.ISOFormat() == "" { // String doesn't contain a date
				times = Times{Unix: nil, Natural: nil}
			} else { // String contains a date.
				day := str.Day()
				month := str.Month()
				year := str.Year()
				date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
				unix = date.Unix()
				natural = date.Format("January 2, 2006")
				times = Times{Unix: unix, Natural: natural}
			}
		}
		js, err := json.Marshal(times)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func main() {
	http.HandleFunc("/", Router)
	http.ListenAndServe(":8080", nil)
}

func Start(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Timestamp Microservice")
}
