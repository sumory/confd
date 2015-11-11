package file

import (
	"strings"

	"fmt"
	"github.com/BurntSushi/toml"
)

type FileStore struct {
	Data map[string]interface{} `toml:"data"`
}

type Client struct {
	FileStore *FileStore
}

func NewFileClient(filePath string) (*Client, error) {
	var fileStore *FileStore
	_, err := toml.DecodeFile(filePath, &fileStore)
	fmt.Printf("%s %v\n", filePath, fileStore)

	if err != nil {
		return nil, fmt.Errorf("Cannot process file store %s - %s", filePath, err.Error())
	}

	return &Client{
		FileStore: fileStore,
	}, nil
}

func (c *Client) GetValues(keys []string) (map[string]interface{}, error) {
	vars := make(map[string]interface{})
	for _, key := range keys {
		key = strings.Trim(key, " ")
		value, ok := c.FileStore.Data[key]
		if ok {
			vars[key] = value
			continue
		}
	}
	return vars, nil
}

func (c *Client) WatchPrefix(prefix string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	<-stopChan
	return 0, nil
}
