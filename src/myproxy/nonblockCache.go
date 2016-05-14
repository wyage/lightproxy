package myproxy

import (
	"errors"
	"fmt"
	"sync"
)

var ErrNotSuitable = errors.New("not suitable for cache")

type entry struct {
	res     result
	ready   chan struct{} //closed when res is ready
	invalid bool          //is it really suitable for cache (in case has Set-Cookie header)
}

type Memo struct {
	requestFunc Func
	cache       map[string]*entry
	mu          sync.Mutex
}

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func NewCache(f Func) *Memo {
	return &Memo{requestFunc: f, cache: make(map[string]*entry)}
}

//unsafe implementation
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		fmt.Println("entry not found:", key)
		e = &entry{ready: make(chan struct{})}
		memo.cache[key] = e
		memo.mu.Unlock()

		fmt.Println("not found", key)

		tres, terr := memo.requestFunc(key)
		if terr != nil {
			if terr == ErrNotSuitable {
				e.invalid = true
				//we see from header that this url is not suitable for cache
				//so we do not update this response to corresponding cache entry
				return tres, nil
			} else {
				//error when requesting
				//todo: consider retrying
				e.res.value, e.res.err = nil, terr
			}
		} else {
			//save response to cache entry
			e.res.value, e.res.err = tres, nil
		}

		close(e.ready)
	} else if e.invalid {
		memo.mu.Unlock()
		//this url is not suitable for cache, we know this from the real response
		return nil, ErrNotSuitable
	} else {
		memo.mu.Unlock()
		<-e.ready
	}

	return e.res.value, e.res.err
}
