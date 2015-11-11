package cli

import (
	"errors"
	"github.com/sumory/confd/cli/file"
//	"github.com/sumory/confd/cli/redis"
//	"github.com/sumory/confd/cli/zookeeper"
)

type Cli interface {
	GetAll() (map[string]interface{}, error)
	GetValues(keys []string) (map[string]interface{}, error)
	GetValue(key string) (interface{}, error)
	SetValue(key string, value interface{}) error
	DeleteKey(key string) error
}

func New(c *CliConfig) (error, Cli) {
	if c.Store == "" {
		c.Store = "file" //default cli
	}

	connectAddr := c.ConnectAddr
	switch c.Store {
	//	case "zookeeper":
	//		return zookeeper.NewZookeeperCli(connectAddr)
	//	case "redis":
	//		return redis.NewRedisCli(connectAddr)
	case "file":
		return file.NewFileCli(connectAddr)
	}

	return errors.New("Invalid cli..."), nil
}
