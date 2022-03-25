package file_cache

import (
	"os"
	"sync"
)

type Cache struct {
	mu       sync.Mutex
	capacity int
	queue    *List
	items    listItems
}

type listItems map[string]*ListItem

type cachedImage struct {
	url  string
	file *os.File
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(listItems, capacity),
	}
}

func (c *Cache) Remember(key string, callback func() *os.File) *os.File {
	fromCache := c.get(key)
	if fromCache != nil {
		return fromCache
	}

	file := callback()
	if ok := c.set(key, file); !ok {
		return nil
	}

	return file
}

func (c *Cache) set(key string, file *os.File) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	newItem := &cachedImage{
		url:  key,
		file: file,
	}
	if _, ok := c.items[key]; ok { // item exists
		c.items[key].Value = newItem
		c.queue.MoveToFront(c.items[key])

		return true
	}

	listItem := c.queue.PushFront(newItem)
	c.items[key] = listItem
	if c.queue.Len() > c.capacity {
		lastElem := c.queue.Back()
		c.queue.Remove(lastElem)
		delete(c.items, lastElem.Value.(*cachedImage).url)
		// todo remove item from disk
	}

	return false
}

func (c *Cache) get(key string) *os.File {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return nil
	}
	c.queue.MoveToFront(elem)

	return elem.Value.(*cachedImage).file
}

func (c *Cache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = listItems{}
}
