package http

import (
	"context"
	"net/http"
	"testing"

	mock_http "github.com/arthurshafikov/image-previewer/internal/transport/http/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestServe(t *testing.T) {
	ctrl := gomock.NewController(t)
	handlerMock := mock_http.NewMockHandler(ctrl)
	loggerMock := mock_http.NewMockLogger(ctrl)
	server := NewServer(loggerMock, handlerMock)
	ctx, cancel := context.WithCancel(context.Background())
	group, ctx := errgroup.WithContext(ctx)
	gomock.InOrder(
		handlerMock.EXPECT().Init(gomock.Any()),
		loggerMock.EXPECT().Info(gomock.Any()),
	)
	group.Go(func() error {
		defer cancel()

		response, err := http.Get("http://localhost:9999") //nolint:noctx
		require.NoError(t, response.Body.Close())
		require.NoError(t, err)

		return nil
	})

	server.Serve(ctx, group, "9999")

	require.NoError(t, group.Wait())
}
