package main

import (
	"io/ioutil"
	"net/http"
	"sync"
)

type SimpleCache struct {
	sync.RWMutex

	cacheEntry *SimpleCacheEntry
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{}
}

func (sc *SimpleCache) CacheEntryExist() (cacheExist bool) {
	sc.RLock()
	cacheExist = sc.cacheEntry != nil
	sc.RUnlock()

	return cacheExist
}

func (sc *SimpleCache) GetCacheEntry() (contentType string, bytes []byte, err error) {
	sc.RLock()
	contentType, bytes, err = sc.cacheEntry.Get()
	sc.RUnlock()

	return contentType, bytes, err
}

func (sc *SimpleCache) UpdateCacheEntry(contentType string, bytes []byte) (err error) {
	sc.RLock()
	err = sc.cacheEntry.Set(contentType, bytes)
	sc.RUnlock()

	return err
}

func (sc *SimpleCache) CreateCacheEntry(contentType string, bytes []byte) {
	sc.Lock()
	sc.cacheEntry = NewSimpleCacheEntry(contentType, bytes)
	sc.Unlock()

	return
}

func (sc *SimpleCache) DeleteCacheEntry() {
	sc.Lock()
	sc.cacheEntry = nil
	sc.Unlock()

	return
}

func (sc *SimpleCache) HttpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		sc.HttpGetHandler(w, r)
	case "POST":
		sc.HttpPostHandler(w, r)
	case "PUT":
		sc.HttpPutHandler(w, r)
	case "DELETE":
		sc.HttpDeleteHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (sc *SimpleCache) HttpGetHandler(w http.ResponseWriter, r *http.Request) {
	if !sc.CacheEntryExist() {
		http.Error(w, "cacheEntry not exists", http.StatusNotFound)
		return
	}

	contentType, bytes, err := sc.GetCacheEntry()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (sc *SimpleCache) HttpPostHandler(w http.ResponseWriter, r *http.Request) {
	if sc.CacheEntryExist() {
		http.Error(w, "cacheEntry already exists", http.StatusConflict)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentType := r.Header.Get("Content-Type")

	sc.CreateCacheEntry(contentType, bytes)

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusCreated)
	w.Write(bytes)
}

func (sc *SimpleCache) HttpPutHandler(w http.ResponseWriter, r *http.Request) {
	if !sc.CacheEntryExist() {
		http.Error(w, "cacheEntry not exists", http.StatusNotFound)
		return
	}

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

	err = sc.UpdateCacheEntry(contentType, bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (sc *SimpleCache) HttpDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if !sc.CacheEntryExist() {
		http.Error(w, "cacheEntry not exists", http.StatusNotFound)
		return
	}

	sc.DeleteCacheEntry()

	w.WriteHeader(http.StatusOK)
}
