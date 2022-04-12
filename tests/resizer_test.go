package tests

import (
	"io/ioutil"
	"net/http"
	"time"
)

func (s *APITestSuite) TestResizeRemoteHostNotExists() {
	statusCode, body := s.postRequest("/resize/200/100/" + remoteHostNotExists)
	r.Equal(http.StatusUnprocessableEntity, statusCode)
	r.Equal(errorHostNotFound, body)
}

func (s *APITestSuite) TestResizeRemoteHostReturns404() {
	statusCode, body := s.postRequest("/resize/200/100/" + remoteHostImageURLNotExists)
	r.Equal(http.StatusUnprocessableEntity, statusCode)
	r.Equal(error404Response, body)
}

func (s *APITestSuite) TestResizeRemoteHostReturnsExeFile() {
	statusCode, body := s.postRequest("/resize/200/100/" + remoteHostImageURLExeFile)
	r.Equal(http.StatusUnprocessableEntity, statusCode)
	r.Equal(errorNotJpgJpegFileResponse, body)
}

func (s *APITestSuite) TestResizeAExceedNumberOfResizedImages() {
	statusCode, _ := s.postRequest("/resize/300/200/" + remoteHostImage1URL)
	r.Equal(http.StatusOK, statusCode)
	statusCode, _ = s.postRequest("/resize/200/200/" + remoteHostImage1URL)
	r.Equal(http.StatusOK, statusCode)
	statusCode, _ = s.postRequest("/resize/100/200/" + remoteHostImage1URL)
	r.Equal(http.StatusOK, statusCode)
	statusCode, _ = s.postRequest("/resize/50/200/" + remoteHostImage1URL)
	r.Equal(http.StatusOK, statusCode)
	statusCode, _ = s.postRequest("/resize/50/100/" + remoteHostImage1URL)
	r.Equal(http.StatusOK, statusCode)

	storageRawFolder, err := ioutil.ReadDir("./storage/raw")
	r.NoError(err)
	rawImages := []string{}
	for _, v := range storageRawFolder {
		rawImages = append(rawImages, v.Name())
	}
	storageResizedFolder, err := ioutil.ReadDir("./storage/resized")
	r.NoError(err)
	resizedImages := []string{}
	for _, v := range storageResizedFolder {
		resizedImages = append(resizedImages, v.Name())
	}

	expectedRawImages := []string{"test-image-1.jpg"}
	expectedResizedImages := []string{"test-image-1_100x200.jpg", "test-image-1_50x200.jpg", "test-image-1_50x100.jpg"}
	r.ElementsMatch(expectedRawImages, rawImages)
	r.ElementsMatch(expectedResizedImages, resizedImages)
}

func (s *APITestSuite) TestResizeDefault() {
	statusCode, _ := s.postRequest("/resize/300/200/" + remoteHostImage1URL)
	r.Equal(http.StatusOK, statusCode)
	statusCode, _ = s.postRequest("/resize/500/200/" + remoteHostImage2URL)
	r.Equal(http.StatusOK, statusCode)
	statusCode, _ = s.postRequest("/resize/100/100/" + remoteHostImage3URL)
	r.Equal(http.StatusOK, statusCode)

	storageRawFolder, err := ioutil.ReadDir("./storage/raw")
	r.NoError(err)
	rawImages := []string{}
	for _, v := range storageRawFolder {
		rawImages = append(rawImages, v.Name())
	}
	storageResizedFolder, err := ioutil.ReadDir("./storage/resized")
	r.NoError(err)
	resizedImages := []string{}
	for _, v := range storageResizedFolder {
		resizedImages = append(resizedImages, v.Name())
	}

	expectedRawImages := []string{"test-image-1.jpg", "test-image-2.jpeg", "test-image-3.jpg"}
	expectedResizedImages := []string{"test-image-1_300x200.jpg", "test-image-2_500x200.jpeg", "test-image-3_100x100.jpg"}
	r.ElementsMatch(expectedRawImages, rawImages)
	r.ElementsMatch(expectedResizedImages, resizedImages)
}

func (s *APITestSuite) TestResizedCheckThatCachedImageIsUsed() {
	startTime := time.Now()
	statusCode, _ := s.postRequest("/resize/300/200/" + remoteHostImage1URL)
	endTime := time.Now()
	r.Equal(http.StatusOK, statusCode)
	requestDurationNothingCached := endTime.Sub(startTime)

	startTime = time.Now()
	statusCode, _ = s.postRequest("/resize/500/100/" + remoteHostImage1URL)
	endTime = time.Now()
	r.Equal(http.StatusOK, statusCode)
	requestDurationRawImageIsCached := endTime.Sub(startTime)

	startTime = time.Now()
	statusCode, _ = s.postRequest("/resize/300/200/" + remoteHostImage1URL)
	endTime = time.Now()
	r.Equal(http.StatusOK, statusCode)
	requestDurationResizedImageIsCached := endTime.Sub(startTime)

	r.True(requestDurationNothingCached > requestDurationRawImageIsCached*5)
	r.True(requestDurationNothingCached > requestDurationResizedImageIsCached*10)
}
