package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type Limit struct {
	count int
	start time.Time
}

var users = map[string]*Limit{}

func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		u := users[ip]
		if u == nil || time.Since(u.start) >= time.Minute {
			u = &Limit{0, time.Now()}
			users[ip] = u
		}

		u.count++

		if u.count > 5 {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/", rateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello!")
	})))

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}