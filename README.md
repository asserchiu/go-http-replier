# go-http-replier
Simple HTTP server with some useful tools

## Routes

* / => HTTP 200
* /hostname => HTTP 200 `os.Hostname()`
* /uptime => HTTP 200 `uptime.String()`
* /ping => HTTP 200 `"pong"`
* /healthz => HTTP 200 `"ok"`
* /healthz10s
    * uptime less equal than 10s => HTTP 200 `"ok"`
    * uptime over 10s => HTTP 500 `fmt.Sprintf("error: %v", uptime.String())`
* /echo =>
    * response with request body
* /echoAll =>
    * JSON encoded response with request body, header, form ...
* /simpleCache
    * POST =>
        * Set request body to cacheEntry.
        * Failed when cacheEntry exist.
    * PUT =>
        * Set request body to cacheEntry.
        * Failed when cacheEntry not exist.
    * GET =>
        * Get data from cacheEntry.
        * Failed when cacheEntry not exist.
    * DELETE =>
        * Delete data from cacheEntry.
        * Failed when cacheEntry not exist.
* /etag
    * .html => a web page with browser sent or new ETag value inside
    * .js => a executable script that can set window.ETag as browser sent or new ETag value
    * => Text string of browser sent or new ETag value

## ENV

* ADDR
    * default: `":8080"`
