package image_cache //nolint

import (
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/arthurshafikov/image-previewer/internal/core"
	"github.com/stretchr/testify/require"
)

var (
	someImageURL  = "https://some-site.com/some-image.jpg"
	someImage     = &core.Image{}
	someImageURL2 = "https://some-site2.com/123some-image.jpg"
	someImage2    = &core.Image{}
	someImageURL3 = "https://some-site3.com/333some-image.jpg"
	someImage3    = &core.Image{}
	someImageURL4 = "https://some-site4.com/444some-image.jpg"
	someImage4    = &core.Image{}
)

func TestEmptyCache(t *testing.T) {
	c := NewCache(10, "")

	image := c.get(someImageURL)
	require.Nil(t, image)
	image = c.get(someImageURL2)
	require.Nil(t, image)
}

func TestSimple(t *testing.T) {
	c := NewCache(5, "")

	deletedImage, err := c.set(someImageURL, someImage)
	require.NoError(t, err)
	require.Nil(t, deletedImage)
	deletedImage, err = c.set(someImageURL2, someImage2)
	require.NoError(t, err)
	require.Nil(t, deletedImage)

	require.Equal(t, someImage, c.get(someImageURL))
	require.Equal(t, someImage2, c.get(someImageURL2))
}

func TestLowCapacity(t *testing.T) {
	c := NewCache(1, "")

	deletedImage, err := c.set(someImageURL, someImage)
	require.NoError(t, err)
	require.Nil(t, deletedImage)
	deletedImage, err = c.set(someImageURL2, someImage2)
	require.NoError(t, err)
	require.Equal(t, someImage, deletedImage)
	deletedImage, err = c.set(someImageURL3, someImage3)
	require.NoError(t, err)
	require.Equal(t, someImage2, deletedImage)

	require.Nil(t, c.get(someImageURL))
	require.Nil(t, c.get(someImageURL2))
	require.Equal(t, someImage3, c.get(someImageURL3))
}

func TestClear(t *testing.T) {
	c := NewCache(2, "")

	deletedImage, err := c.set(someImageURL, someImage)
	require.NoError(t, err)
	require.Nil(t, deletedImage)
	deletedImage, err = c.set(someImageURL2, someImage2)
	require.NoError(t, err)
	require.Nil(t, deletedImage)

	c.Clear()

	result := c.get(someImageURL)
	require.Nil(t, result)

	result = c.get(someImageURL2)
	require.Nil(t, result)
}

func TestUnknownKeys(t *testing.T) {
	c := NewCache(5, "")

	result := c.get(someImageURL)
	require.Nil(t, result)
	result = c.get(someImageURL2)
	require.Nil(t, result)
}

func TestPurgeOldElement(t *testing.T) {
	c := NewCache(3, "")

	deletedImage, err := c.set(someImageURL, someImage)
	require.NoError(t, err)
	require.Nil(t, deletedImage)
	deletedImage, err = c.set(someImageURL2, someImage2)
	require.NoError(t, err)
	require.Nil(t, deletedImage)
	deletedImage, err = c.set(someImageURL3, someImage3)
	require.NoError(t, err)
	require.Nil(t, deletedImage)

	deletedImage, err = c.set(someImageURL2, someImage2)
	require.NoError(t, err)
	require.Nil(t, deletedImage)

	c.get(someImageURL2)
	c.get(someImageURL)

	deletedImage, err = c.set(someImageURL2, someImage2)
	require.NoError(t, err)
	require.Nil(t, deletedImage)

	deletedImage, err = c.set(someImageURL2, someImage2)
	require.NoError(t, err)
	require.Nil(t, deletedImage)

	deletedImage, err = c.set(someImageURL4, someImage4)
	require.NoError(t, err)
	require.Equal(t, someImage3, deletedImage)

	result := c.get(someImageURL3)
	require.Nil(t, result)
}

func TestRemember(t *testing.T) {
	c := NewCache(1, "")

	imageFile, err := os.Create("img.jpg")
	require.NoError(t, err)
	defer imageFile.Close()

	image := &core.Image{
		Name:      "img",
		Extension: "jpg",
		File:      imageFile,
	}
	require.NoError(t, imageFile.Close())

	result, err := c.Remember(someImageURL, func() (*core.Image, error) {
		return image, nil
	})
	require.Equal(t, image, result)
	require.NoError(t, err)

	result, err = c.Remember(someImageURL2, func() (*core.Image, error) {
		return someImage, nil
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
			_, err := c.set(strconv.Itoa(i), someImage)
			require.NoError(t, err)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.get(strconv.Itoa(rand.Intn(1_000_000))) //nolint:gosec
		}
	}()

	wg.Wait()
}
