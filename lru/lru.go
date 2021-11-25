// lru:Least Recently Used
package lru

import "container/list"

// 带自动删除lru的Cache
type Cache struct {
	// cache最大容量,超过最大容量触发删除lru
	MaxEntries int
	// 删除lru触发的回调
	OnEvicted func(key Key, value interface{})

	// list与map共同缓存了所有数据,list增删快,map查询快(go里面map是散列表)
	ll    *list.List
	cache map[Key]*list.Element
}

type Key interface{}

// list.Element.Value中存的 真实的value
type entry struct {
	key   Key
	value interface{}
}

// 获取Cache实例,如果macEntries为0,则不会有lru功能,需要使用者手动删除缓存数据
func New(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[Key]*list.Element),
		// 这里为什么没有初始化OnEvicted,有点奇怪!!
	}
}

func (c *Cache) Add(key Key, value interface{}) {
	// 所有的操作函数前面为什么都要判断一下?
	// 是因为下面的Clear()函数 c.cache=nil   c.ll=nil
	if c.cache == nil {
		c.cache = make(map[Key]*list.Element)
		c.ll = list.New()
	}
	// 当前key已经缓存了,把当前数据移到链表头部
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}
	// 当前key没有缓存,链表头部插入数据,检查长度,判断要不要删除lru
	ele := c.ll.PushFront(&entry{key: key, value: value})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

// 获取key对应的缓存数据
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// 删除指定key的数据
func (c *Cache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// 从链表尾部删除一个数据
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

// 删除数据,有回调的话触发回调
func (c *Cache) removeElement(ele *list.Element) {
	c.ll.Remove(ele)
	kv := ele.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// 获取缓存数据的数量(链表的长度就是数据的数量)
func (c *Cache) Len() int {
	return c.ll.Len()
}

// 清空缓存数据
func (c *Cache) Clear() {
	c.ll = nil
	c.cache = nil
}
