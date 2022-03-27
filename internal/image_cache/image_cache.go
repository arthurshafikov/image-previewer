package image_cache

import (
	"os"
	"sync"

	"github.com/thewolf27/image-previewer/internal/core"
)

type Cache struct {
	mu                 sync.Mutex
	capacity           int
	queue              *List
	items              listItems
	cachedImagesFolder string
}

type listItems map[string]*ListItem

type cachedImage struct {
	url   string
	image *core.Image
}

func NewCache(capacity int, cachedImagesFolder string) *Cache {
	return &Cache{
		capacity:           capacity,
		queue:              NewList(),
		items:              make(listItems, capacity),
		cachedImagesFolder: cachedImagesFolder,
	}
}

func (c *Cache) GetCachedImagesFolder() string {
	return c.cachedImagesFolder
}

func (c *Cache) Remember(key string, callback func() (*core.Image, error)) (*core.Image, error) {
	fromCache := c.get(key)
	if fromCache != nil {
		return fromCache, nil
	}

	image, err := callback()
	if err != nil {
		return nil, err
	}
	deletedImage, err := c.set(key, image)
	if err != nil {
		return nil, err
	}
	if deletedImage != nil {
		if err := os.Remove(deletedImage.File.Name()); err != nil {
			return nil, err
		}
	}

	return image, nil
}

func (c *Cache) set(key string, image *core.Image) (deletedImage *core.Image, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	newItem := &cachedImage{
		url:   key,
		image: image,
	}
	if _, ok := c.items[key]; ok { // item exists
		c.items[key].Value = newItem
		c.queue.MoveToFront(c.items[key])

		return nil, nil
	}

	listItem := c.queue.PushFront(newItem)
	c.items[key] = listItem
	if c.queue.Len() > c.capacity {
		lastElem := c.queue.Back()
		c.queue.Remove(lastElem)
		delete(c.items, lastElem.Value.(*cachedImage).url)
		return lastElem.Value.(*cachedImage).image, nil
	}

	return nil, nil
}

func (c *Cache) get(key string) *core.Image {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return nil
	}
	c.queue.MoveToFront(elem)

	return elem.Value.(*cachedImage).image
}

func (c *Cache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = listItems{}
}
