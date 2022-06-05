package singleflight

import "sync"

// call 代表正在进行中或已经结束的请求,通过waitgroup来防止重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group 管理不同key的请求
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// Do 相同的key，不管Do被调用多少次，fn只执行一次
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
