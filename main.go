package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	startTime := time.Now()

	http.HandleFunc("/uptime", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Now().Sub(startTime)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(uptime.String()))
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/healthz10s", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Now().Sub(startTime)
		if uptime.Seconds() > 10 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error: %v", uptime.String())))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}
	})

	http.HandleFunc("/echo", echoRequestHandler)
	http.HandleFunc("/echoAll", echoAllRequestHandler)

	http.HandleFunc("/simpleCache", NewSimpleCache().HttpHandler)

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Fatal(http.ListenAndServe(addr, nil))
}

func detectContentType(bytes []byte) (contentType string) {
	err := json.Unmarshal(bytes, &struct{}{})
	if err == nil {
		return "application/json"
	}
	return http.DetectContentType(bytes)
}

func echoRequestHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = detectContentType(bytes)
	}

	w.Header().Set("Content-Type", contentType)

	w.Write(bytes)
}

func echoAllRequestHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(map[string]interface{}{
		"Body":                   string(bytes),
		"ContentLength":          r.ContentLength,
		"Form":                   r.Form,
		"Header":                 r.Header,
		"Host":                   r.Host,
		"Method":                 r.Method,
		"Proto":                  r.Proto,
		"URL":                    r.URL,
		"URL.Query":              r.URL.Query(),
		"http.DetectContentType": http.DetectContentType(bytes),
		"json.Unmarshal error":   json.Unmarshal(bytes, &struct{}{}),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
