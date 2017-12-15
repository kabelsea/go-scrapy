package scrapy

import "sync"

type Context struct {
	contextMap map[string]interface{}
	lock       *sync.RWMutex
}

// Create Request or Response context instance
func NewContext() *Context {
	return &Context{
		contextMap: make(map[string]interface{}),
		lock:       &sync.RWMutex{},
	}
}
