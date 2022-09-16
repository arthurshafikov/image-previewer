package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFullName(t *testing.T) {
	image := Image{
		Name:      "someName",
		Extension: "ext",
	}

	require.Equal(t, "someName.ext", image.GetFullName())
}

func TestGetFullNameWithWidthAndHeight(t *testing.T) {
	testCases := []struct {
		Image    Image
		Width    int
		Height   int
		Expected string
	}{
		{
			Image: Image{
				Name:      "name1",
				Extension: "jpg",
			},
			Width:    100,
			Height:   200,
			Expected: "name1_100x200.jpg",
		},
		{
			Image: Image{
				Name:      "name2",
				Extension: "png",
			},
			Width:    500,
			Height:   200,
			Expected: "name2_500x200.png",
		},
		{
			Image: Image{
				Name:      "name3",
				Extension: "gif",
			},
			Width:    10,
			Height:   50,
			Expected: "name3_10x50.gif",
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.Expected, testCase.Image.GetFullNameWithWidthAndHeight(testCase.Width, testCase.Height))
	}
}
