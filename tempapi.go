package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var AccessToken string

func TempPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		q := r.URL.Query()

		// Auth
		token := q.Get("token")
		if token == "" {
			http.Error(w, "No token", http.StatusUnauthorized)
			return
		}
		if token != AccessToken {
			http.Error(w, "Wrong token", http.StatusForbidden)
			return
		}

		// Get time
		timestring := q.Get("time")
		if timestring == "" {
			http.Error(w, "No time given", http.StatusBadRequest)
			return
		}
		timeint, err := strconv.ParseInt(timestring, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		timeval := time.Unix(timeint, 0)
		fmt.Println(timeval)

		// Get room
		room := q.Get("room")
		if room == "" {
			http.Error(w, "No room given", http.StatusBadRequest)
			return
		}

		// Get temp
		tempstring := q.Get("temp")
		if tempstring == "" {
			http.Error(w, "No temp given", http.StatusBadRequest)
			return
		}
		temp, err := strconv.ParseFloat(tempstring, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("Room:", room)
		fmt.Println("Time:", timeval)
		fmt.Println("Temperature:", float32(temp))
	} else {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	TLSKey := os.Getenv("tlskey")
	TLSCert := os.Getenv("tlscert")
	AccessToken = os.Getenv("token")
	http.HandleFunc("/", TempPost)
	log.Fatal(http.ListenAndServeTLS(":8443", TLSCert, TLSKey, nil))
}
