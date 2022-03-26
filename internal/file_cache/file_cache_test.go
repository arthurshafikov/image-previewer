package file_cache

import (
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thewolf27/image-previewer/internal/core"
)

var (
	someImageUrl  = "https://some-site.com/some-image.jpg"
	someImage     = &core.Image{}
	someImageUrl2 = "https://some-site2.com/123some-image.jpg"
	someImage2    = &core.Image{}
	someImageUrl3 = "https://some-site3.com/333some-image.jpg"
	someImage3    = &core.Image{}
	someImageUrl4 = "https://some-site4.com/444some-image.jpg"
	someImage4    = &core.Image{}
)

func TestEmptyCache(t *testing.T) {
	c := NewCache(10, "")

	image := c.get(someImageUrl)
	require.Nil(t, image)
	image = c.get(someImageUrl2)
	require.Nil(t, image)
}

func TestSimple(t *testing.T) {
	c := NewCache(5, "")

	deletedImage, err := c.set(someImageUrl, someImage)
	require.Nil(t, deletedImage)
	require.NoError(t, err)
	deletedImage, err = c.set(someImageUrl2, someImage2)
	require.Nil(t, deletedImage)
	require.NoError(t, err)

	require.Equal(t, someImage, c.get(someImageUrl))
	require.Equal(t, someImage2, c.get(someImageUrl2))
}

func TestLowCapacity(t *testing.T) {
	c := NewCache(1, "")

	deletedImage, err := c.set(someImageUrl, someImage)
	require.Nil(t, deletedImage)
	require.NoError(t, err)
	deletedImage, err = c.set(someImageUrl2, someImage2)
	require.Equal(t, someImage, deletedImage)
	require.NoError(t, err)
	deletedImage, err = c.set(someImageUrl3, someImage3)
	require.Equal(t, someImage2, deletedImage)
	require.NoError(t, err)

	require.Nil(t, c.get(someImageUrl))
	require.Nil(t, c.get(someImageUrl2))
	require.Equal(t, someImage3, c.get(someImageUrl3))
}

func TestClear(t *testing.T) {
	c := NewCache(2, "")

	deletedImage, err := c.set(someImageUrl, someImage)
	require.Nil(t, deletedImage)
	require.NoError(t, err)
	deletedImage, err = c.set(someImageUrl2, someImage2)
	require.Nil(t, deletedImage)
	require.NoError(t, err)

	c.clear()

	result := c.get(someImageUrl)
	require.Nil(t, result)

	result = c.get(someImageUrl2)
	require.Nil(t, result)
}

func TestUnknownKeys(t *testing.T) {
	c := NewCache(5, "")

	result := c.get(someImageUrl)
	require.Nil(t, result)
	result = c.get(someImageUrl2)
	require.Nil(t, result)
}

func TestPurgeOldElement(t *testing.T) {
	c := NewCache(3, "")

	deletedImage, err := c.set(someImageUrl, someImage)
	require.Nil(t, deletedImage)
	require.NoError(t, err)
	deletedImage, err = c.set(someImageUrl2, someImage2)
	require.Nil(t, deletedImage)
	require.NoError(t, err)
	deletedImage, err = c.set(someImageUrl3, someImage3)
	require.Nil(t, deletedImage)
	require.NoError(t, err)

	deletedImage, err = c.set(someImageUrl2, someImage2)
	require.Nil(t, deletedImage)
	require.NoError(t, err)

	c.get(someImageUrl2)
	c.get(someImageUrl)

	deletedImage, err = c.set(someImageUrl2, someImage2)
	require.Nil(t, deletedImage)
	require.NoError(t, err)

	deletedImage, err = c.set(someImageUrl2, someImage2)
	require.Nil(t, deletedImage)
	require.NoError(t, err)

	deletedImage, err = c.set(someImageUrl4, someImage4)
	require.Equal(t, someImage3, deletedImage)
	require.NoError(t, err)

	result := c.get(someImageUrl3)
	require.Nil(t, result)
}

func TestRemember(t *testing.T) {
	c := NewCache(1, "")

	imageFile, err := os.Create("img.jpg")
	require.NoError(t, err)

	image := &core.Image{
		Name:      "img",
		Extension: "jpg",
		File:      imageFile,
	}
	require.NoError(t, imageFile.Close())

	result, err := c.Remember(someImageUrl, func() *core.Image {
		return image
	})
	require.Equal(t, image, result)
	require.NoError(t, err)

	result, err = c.Remember(someImageUrl2, func() *core.Image {
		return someImage
	})
	require.Equal(t, someImage, result)
	require.NoError(t, err)

	_, err = os.Open(imageFile.Name())
	require.NotNil(t, err)
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10, "")
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.set(strconv.Itoa(i), someImage)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.get(strconv.Itoa(rand.Intn(1_000_000)))
		}
	}()

	wg.Wait()
}
