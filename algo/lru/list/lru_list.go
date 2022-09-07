package lru

import (
	"sync"
)

// Node 存储数据结构体
type Node struct {
	key   interface{}
	value interface{}
	prev  *Node // 往 first 方向
	next  *Node // 往 last 方向
}

// Cache 实现了LRU的结构体
type Cache struct {
	mux   sync.RWMutex
	len   int   // 当前长度
	cap   int   // 最大容量
	first *Node // 队首（最右边），最常使用的
	last  *Node // 队尾（最左边），最少使用的
	nodes map[interface{}]*Node
}

// NewLRUCache 并发安全的LRU缓存
func NewLRUCache(capacity int) *Cache {
	first := &Node{}
	last := &Node{prev: first}
	first.next = last

	return &Cache{
		len:   0,
		cap:   capacity,
		first: first,
		last:  last,
		nodes: make(map[interface{}]*Node, capacity),
	}
}

func (l *Cache) Get(key interface{}) interface{} {
	l.mux.RLock()
	defer l.mux.RUnlock()

	if node, ok := l.nodes[key]; ok {
		// 将该值设为第一个
		l.moveToFirst(node)
		return node.value
	}
	return nil
}

func (l *Cache) Put(key, value interface{}) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if node, ok := l.nodes[key]; ok {
		// 更新值
		node.value = value
		l.moveToFirst(node)
	} else {
		// 到达最大容量了，删除最后面的值
		if l.len == l.cap {
			delete(l.nodes, l.last.prev.key)
			l.removeLast()
		} else {
			l.len++
		}
	}
	node := &Node{
		key:   key,
		value: value,
	}
	l.nodes[key] = node
	l.insertToFirst(node)
}

func (l *Cache) moveToFirst(node *Node) {
	// 1. 将该节点从链表里面删掉
	switch node {
	case l.first.next:
		// 队首，不做改变
		return
	default:
		// 队中，删掉该节点
		node.prev.next = node.next
		node.next.prev = node.prev
	}
	// 2. 将该节点插入到队首
	l.insertToFirst(node)
}

func (l *Cache) removeLast() {
	if l.last.prev == l.first {
		return
	}
	// 双向链表长度大于1
	l.last.prev = l.last.prev.prev
	l.last.prev.next = l.last
}

func (l *Cache) insertToFirst(node *Node) {
	l.first.next.prev = node
	node.next = l.first.next
	node.prev = l.first
	l.first.next = node
}

func (l *Cache) Keys() []interface{} {
	var keys []interface{}
	var cur = l.last.prev
	for cur != l.first {
		keys = append(keys, cur.key)
		cur = cur.prev
	}
	return keys
}
