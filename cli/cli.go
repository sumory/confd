package cli

import (
	"errors"
	"github.com/sumory/confd/cli/file"
	//	"github.com/sumory/confd/cli/redis"
	//	"github.com/sumory/confd/cli/zookeeper"
)

type Cli interface {
	GetValues(keys []string) (map[string]interface{}, error)
	SetValues(values map[string]interface{}) error
}

func New(c *CliConfig) (Cli, error) {
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

	return nil, errors.New("Invalid cli...")
}
