package store

import (
	"errors"
	"github.com/sumory/confd/store/file"
	"github.com/sumory/confd/store/redis"
	"github.com/sumory/confd/store/zookeeper"
)

type StoreClient interface {
	GetValues(keys []string) (map[string]interface{}, error)
	WatchPrefix(prefix string, waitIndex uint64, stopChan chan bool) (uint64, error)
}

func New(c *StoreConfig) (StoreClient, error) {
	if c.Store == "" {
		c.Store = "file" //default store
	}

	connectAddr := c.ConnectAddr
	switch c.Store {
	case "zookeeper":
		return zookeeper.NewZookeeperClient(connectAddr)
	case "redis":
		return redis.NewRedisClient(connectAddr)
	case "file":
		return file.NewFileClient(connectAddr)
	}

	return nil, errors.New("Invalid store...")
}
