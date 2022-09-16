package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	storageFolder     = "./storage"
	configFolderPath  = "./fakeconfigs"
	configFileContent = []byte(`AppConfig:
    SizeOfLRUCacheForRawImages: 10
    SizeOfLRUCacheForResizedImages: 15
ServerConfig:
    Port: 8125`)

	expectedAppConfig = AppConfig{
		SizeOfLRUCacheForRawImages:     10,
		SizeOfLRUCacheForResizedImages: 15,
	}
	expectedServerConfig = ServerConfig{
		Port: "8125",
	}
	expectedStorageConfig = StorageConfig{
		StorageFolder: storageFolder,
	}
)

func TestNewConfigFromEnvFile(t *testing.T) {
	createFakeConfigFileAndFolder(t)
	defer deleteFakeConfigFileAndFolder(t)

	config := NewConfig(configFolderPath, storageFolder)

	require.Equal(t, expectedAppConfig, config.AppConfig)
	require.Equal(t, expectedServerConfig, config.ServerConfig)
	require.Equal(t, expectedStorageConfig, config.StorageConfig)
}

func createFakeConfigFileAndFolder(t *testing.T) {
	t.Helper()
	if err := os.Mkdir(configFolderPath, 0777); err != nil { //nolint:gofumpt
		t.Fatal(err)
	}
	err := os.WriteFile(fmt.Sprintf("%s/config.yml", configFolderPath), configFileContent, 0600) //nolint:gofumpt
	if err != nil {
		t.Fatal(err)
	}
}

func deleteFakeConfigFileAndFolder(t *testing.T) {
	t.Helper()
	if err := os.RemoveAll(configFolderPath); err != nil {
		t.Fatal(err)
	}
}
