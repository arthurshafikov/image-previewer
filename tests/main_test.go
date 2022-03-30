package tests

import (
	"context"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/thewolf27/image-previewer/internal/config"
	"github.com/thewolf27/image-previewer/internal/image_cache"
	"github.com/thewolf27/image-previewer/internal/services"
	server "github.com/thewolf27/image-previewer/internal/transport/http"
	"github.com/thewolf27/image-previewer/internal/transport/http/handler"
)

var (
	r *require.Assertions
)

type APITestSuite struct {
	suite.Suite

	ServerEngine *gin.Engine

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	r = s.Require()
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	defer s.ctxCancel()

	config := config.NewConfig("./configs", "./storage")

	rawImageCache := image_cache.NewCache(config.AppConfig.SizeOfLRUCacheForRawImages, "./storage/raw")
	resizedImageCache := image_cache.NewCache(config.AppConfig.SizeOfLRUCacheForResizedImages, "./storage/resized")
	services := services.NewServices(services.Deps{
		Config:            config,
		RawImageCache:     rawImageCache,
		ResizedImageCache: resizedImageCache,
	})
	handler := handler.NewHandler(s.ctx, services)
	server := server.NewServer(handler)
	s.ServerEngine = server.Engine

	handler.Init(s.ServerEngine)
}

func (s *APITestSuite) SetupTest() {
	r.NoError(os.MkdirAll("storage/raw", fs.ModePerm))
	r.NoError(os.MkdirAll("storage/resized", fs.ModePerm))
}

func (s *APITestSuite) TearDownTest() {
	r.NoError(os.RemoveAll("storage"))
}

func (s *APITestSuite) TearDownSuite() {
	s.ctxCancel()
}

func (s *APITestSuite) postRequest(route string) (int, string) {
	req := httptest.NewRequest(http.MethodPost, route, nil)
	recorder := httptest.NewRecorder()
	s.ServerEngine.ServeHTTP(recorder, req)
	r.NoError(req.Body.Close())

	recorderResult := recorder.Result()
	bodyBytes, err := ioutil.ReadAll(recorderResult.Body)
	r.NoError(err)
	r.NoError(recorderResult.Body.Close())

	return recorderResult.StatusCode, string(bodyBytes)
}