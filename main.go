package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/hostname", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hostname))
	})

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

	http.HandleFunc("/etag.html", etagRequestHandler)
	http.HandleFunc("/etag.js", etagRequestHandler)
	http.HandleFunc("/etag", etagRequestHandler)

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

func etagRequestHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	etag := r.Header.Get("If-None-Match")
	if etag == "" {
		etag = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	cacheControl := os.Getenv("ETAG_CACHE_CONTROL")
	if cacheControl == "" {
		cacheControl = "max-age=600"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "ETag")
	w.Header().Set("Cache-Control", cacheControl)
	w.Header().Set("ETag", etag)

	// take the extension
	idx := strings.LastIndex(path, ".")
	ext := ""
	if idx > -1 {
		ext = path[idx:]
	}

	switch ext {
	case ".html":
		w.Header().Set("Content-Type", "text/html")
		f := `
<meta name="viewport" content="width=device-width; initial-scale=1.0; maximum-scale=1.0; user-scalable=0;">

<style>
.floating {
  display: table;
  float: right;
  height: 100%%;
  width: 100%%;
  border: 1px solid red;
}
.floating p {
  display: table-cell; 
  vertical-align: middle; 
  text-align: center;
  font-size: calc(1.5*100vw/%d);
}
</style>

<div class="floating"><p><tt>%s</tt></p></div>
`
		fmt.Fprintf(w, f, len(etag), etag)
		return
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
		f := `!function(){window.ETag="%s"}();`
		fmt.Fprintf(w, f, etag)
		return

	}

	w.Write([]byte(etag))
}
