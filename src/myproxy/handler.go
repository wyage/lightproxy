package myproxy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// get requests match certain file types
// without set-cookie
// path without query string, will be cached

// so for a new request not in cache, we should determin if it is cacheable
// from its response header

type MyHandler struct {
	config *ProxyConfig
	mcache *Memo
}

type cacheItem struct {
	header http.Header
	body   []byte
}

func (h *MyHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Printf("new req %#v\n", req)
	fmt.Printf("new req URL %#v\n", req.URL)

	if req.Method == "GET" && req.URL.RawQuery == "" {
		//suitable for cache
		h.getFromCache(resp, req)
	} else {
		h.getFromOrigin(resp, req)
	}
}

func containsQueryString(tpath string) bool {
	idx := strings.Index(tpath, "?")
	if idx > -1 && idx < len(tpath)-1 {
		return true
	}
	return false
}

//only a test method
func serveHello(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	io.WriteString(resp, "Hello world!")
}

//todo: read from cache
func (h *MyHandler) getFromCache(resp http.ResponseWriter, req *http.Request) {
	//h.mcache.Get()
	url := h.config.BackHost + req.URL.Path
	fmt.Println("request url from cache:", url)

	entry, err := h.mcache.Get(url)
	if err == ErrNotSuitable {
		fmt.Println("no suitable url:", url)
		h.getFromOrigin(resp, req)
	} else {
		fmt.Println("cache hit")
		val := entry.(*cacheItem)

		copyHeader(resp.Header(), val.header)
		io.Copy(resp, bytes.NewBuffer(val.body))
	}

}

func (h *MyHandler) getFromOrigin(resp http.ResponseWriter, req *http.Request) {
	//make a new request to backend server
	h.modifyRequest(req)

	tres, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer tres.Body.Close()
	fmt.Printf("header:%#v\n", tres.Header)

	//copy header first
	for k, v := range tres.Header {
		if len(v) < 2 {
			resp.Header().Set(k, v[0])
		} else {
			resp.Header().Set(k, strings.Join(v, ""))
		}
	}
	//copy response body
	io.Copy(resp, tres.Body)
}

//modify new request according configuration
func (h *MyHandler) modifyRequest(req *http.Request) {
	req.URL.Host = h.config.BackHost
	req.URL.Scheme = "http"
	req.Host = h.config.BackHost
	req.RequestURI = ""
}

func NewHandler(config *ProxyConfig) http.Handler {
	return &MyHandler{config, NewCache(doRealRequest)}
}

//todo: for reqs not suitable for cache, do request every time
func doRealRequest(key string) (interface{}, error) {

	if strings.HasPrefix(key, "http://") == false {
		key = "http://" + key
	}

	fmt.Println("will send request:", key)
	resp, err := http.Get(key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	setcookie := resp.Header.Get("Set-Cookie")
	fmt.Printf("set cookie: %s, [%s]\n", key, setcookie)
	if setcookie != "" {
		//although Set-Cookie found, we return this response by the way
		return body, ErrNotSuitable
	}

	//succeeds
	return makeCacheItem(resp.Header, body), nil
}

//returns a pointer
func makeCacheItem(header http.Header, body []byte) *cacheItem {
	newheader := make(map[string][]string)
	copyHeader(newheader, header)
	return &cacheItem{newheader, body}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
