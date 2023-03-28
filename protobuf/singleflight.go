package protobuf

import "sync"

//call代表正在进行中，或已经结束的请求，使用sync.WaitGroup锁避免重入。
//Group是singleflight的主数据结构，管理不同的key的请求(call)

type call struct { //代表了正在进行中，或者已经结束的请求，sync.WaitGroup的作用是避免重入锁
	wg  sync.WaitGroup
	val interface{}
	err error
}
type Group1 struct { //管理不同key的不同请求
	mu sync.Mutex
	m  map[string]*call
}

/**Do方法，接收2个参数，第一个参数是key,第二个参数是一个函数fn。Do的作用就是，针对相同的key,
  ,无论Do被调用多少次，函数fn都只会被调用一次，等fn调用结束了，返回值或错误。
*/

func (g *Group1) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         //如果请求正在进行中，则等待
		return c.val, c.err //请求结束，返回结果
	}
	c := new(call)
	c.wg.Add(1)  //发起请求前加锁
	g.m[key] = c //添加g.m,表明key已经有对应的请求在处理了
	g.mu.Unlock()

	c.val, c.err = fn() //调用fn，发起请求
	c.wg.Done()         //请求结束

	g.mu.Lock()
	delete(g.m, key) //更新g.m
	g.mu.Unlock()

	return c.val, c.err //返回结果
}
