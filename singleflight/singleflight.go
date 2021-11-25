// 单航班:
// 并发情况下,多个goroutine调用同一个函数时,只有一个goroutine真正执行,其他的goroutine等着
// 等执行goroutine执行完,所有goroutine获取返回结果
package singleflight

import (
	"sync"
)

// 在执行或在等待结果的函数对应保存的内容
// value err 对于不同的函数都是可变的,sync.WaitGroup是必须的
type call struct {
	wg    sync.WaitGroup
	value interface{}
	err   error
}

// 单航班实例
type Group struct {
	// 做并发控制
	mx sync.Mutex
	// 记录哪些函数正在执行
	m map[string]*call
}

// 调用执行fn
// 如果多个goroutine调用该函数且key值一样,只有一个goroutine执行,其他goroutine等待结果
// key 函数的标识  fn 需要执行的函数
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mx.Lock()
	// map初始化
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	// 发现有其他goroutine在执行相同函数
	if c, ok := g.m[key]; ok {
		// 先把锁放了
		g.mx.Unlock()
		// 等正在执行函数的goroutine的执行完
		c.wg.Wait()
		// 正在执行函数的goroutine的执行完了,拿着结果直接返回
		return c.value, c.err
	}

	// 没有其他goroutine在执行, 创建一个
	c := new(call)
	// 一个goroutine开始，Add(1)，这里最多只会执行到一次，也就是不会并发调用下面的fn()
	c.wg.Add(1)
	// 记录一下key 对应函数正在执行
	g.m[key] = c
	g.mx.Unlock()

	// 真正执行fn()
	c.value, c.err = fn()
	// 告诉其他goroutine,fu()执行完了,可以获取返回结果了
	c.wg.Done()

	// fn()执行完了,删除key    //个人理解:对fn()这种没有输入参数的函数,返回结果一定的情况,可以缓存数据结果,下次计算的时候直接获取结果,不要再一算一次了
	g.mx.Lock()
	delete(g.m, key)
	g.mx.Unlock()

	return c.value, c.err
}
