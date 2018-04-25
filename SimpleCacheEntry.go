package main

import (
	"errors"
	"sync"
)

var (
	errNilCacheEntry = errors.New("nil cacheEntry")
)

type SimpleCacheEntry struct {
	sync.RWMutex

	contentType string
	bytes       []byte
}

func NewSimpleCacheEntry(contentType string, bytes []byte) (entry *SimpleCacheEntry) {
	if len(contentType) == 0 {
		contentType = detectContentType(bytes)
	}

	return &SimpleCacheEntry{
		contentType: contentType,
		bytes:       bytes,
	}
}

func (c *SimpleCacheEntry) Set(contentType string, bytes []byte) (err error) {
	if c == nil {
		return errNilCacheEntry
	}

	if len(contentType) == 0 {
		contentType = detectContentType(bytes)
	}

	c.Lock()
	c.contentType = contentType
	c.bytes = bytes
	c.Unlock()

	return nil
}

func (c *SimpleCacheEntry) Get() (contentType string, bytes []byte, err error) {
	if c == nil {
		return "", []byte{}, errNilCacheEntry
	}

	c.RLock()
	contentType = c.contentType
	bytes = c.bytes
	c.RUnlock()

	return contentType, bytes, nil
}
