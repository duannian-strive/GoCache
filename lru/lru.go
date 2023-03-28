package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nbytes    int64
	l1        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

//entry是链表的节点的存储类型
type entry struct {
	key   string
	value Value
}

//值类型中接口，让返回存储数据的大小
type Value interface {
	Len() int
}

//实例化Cache方法
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		l1:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//查询功能
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		//将访问到的数据移到队尾
		c.l1.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//删除，即淘汰策略，即移除最近最少访问的节点(队列)
func (c *Cache) RemoveOldest() {
	ele := c.l1.Back()
	if ele != nil {
		//删除链表队头节点数据
		c.l1.Remove(ele)
		kv := ele.Value.(*entry)
		//根据key删除map里面的数据
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

//新增或者修改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.l1.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.l1.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

//看list有多少个数据
func (c *Cache) Len() int {
	return c.l1.Len()
}
