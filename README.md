# Image Previewer

[![Go Report Card](https://goreportcard.com/badge/github.com/arthurshafikov/image-previewer)](https://goreportcard.com/report/github.com/arthurshafikov/image-previewer)
![Tests](https://github.com/arthurshafikov/image-previewer/actions/workflows/tests.yml/badge.svg)
![License](https://img.shields.io/github/license/arthurshafikov/image-previewer)

This is the pet-project that resizes images. It accepts *width*, *height*, and *image url* in query params.
So query would look like this:

`https://resizer-service.com/resize/1080/720/https://any-website.com/some-image.jpg`

Where we need to resize the image from url *`https://any-website.com/some-image.jpg`* to sizes **1080** *width* and **720** *height*.

The microservice downloads image from given **url**, caches it in the storage so that it would't be too many queries to remote host and resizes it to the given dimensions.
Resized image will also be cached so querying the same image **url** with the same dimensions wouldn't make any http requests or actually resizing action, it would just return the cached resized image and that would take much less time than the case when image isn't cached.

For caching system I chose **LRU** (Least Recent Used) algorithm so if image is ofter being requested it would stay in cache for a longer time than an image that has been requested only once.

To show how fast the cache is working we can check the results of one of my test queries, where I have queried the image, then queried it again with other dimensions and then query the same image in the same dimensions

**So the results are:**

|  Alias                                        |                                      |
|:----------------------------------------------|-------------------------------------:|
|   Non Cached Image                            |          **0.022770071s**            |
|   Cached Image But Dimensions Are Different   |          **0.002314989s**            |
|   Cached Image With Same Dimensions           |          **0.000115316s**            |

We can see that getting a *cached image* with the *same dimensions* took only **0.0001** seconds! That's faster than getting *non-cached* image in **~200 times**!

# Commands

## Docker-compose

Run the application
```
make up
```

Down the application
```
make down
```

# Tests

This project is covered with unit-tests and with integration tests.

Integration tests are covering following scenarios:

- if remote_host doesn't exists
- if remote_host returns 404 error
- if remote_host returns exe file instead of image
- if cached storage is being overflowed (that the least used image is being deleted)
- if remote_host returns image (default case)
- if requested image is cached (check that request time is much faster)

## Tests commands 

Run unit tests
```
make test
```

Run integration tests
```
make integration-tests
```

Down all the integration test containers (for reset cache)
```
make reset-integration-tests
```
