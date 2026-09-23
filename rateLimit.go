package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Limit struct {
	count int
	start time.Time
}

var (
	users = map[string]*Limit{}
	mu    sync.Mutex
)

func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		mu.Lock()

		u := users[ip]
		if u == nil || time.Since(u.start) >= time.Minute {
			u = &Limit{0, time.Now()}
			users[ip] = u
		}

		u.count++
		tooMany := u.count > 5

		mu.Unlock()

		if tooMany {
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