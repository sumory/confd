package redis

import (
	"github.com/garyburd/redigo/redis"
	"os"
	"strings"
	"time"
)

type Client struct {
	client redis.Conn
}

func NewRedisClient(address string) (*Client, error) {
	var err error
	var conn redis.Conn
	network := "tcp"
	if _, err = os.Stat(address); err == nil {
		network = "unix"
	}
	conn, err = redis.DialTimeout(network, address, time.Second, time.Second, time.Second)
	if err != nil {
		return nil, err
	}
	return &Client{conn}, nil

}

func (c *Client) GetValues(keys []string) (map[string]interface{}, error) {
	vars := make(map[string]interface{})
	for _, key := range keys {
		key = strings.Trim(key, " ")
		value, err := redis.String(c.client.Do("GET", key))
		if err == nil {
			vars[key] = value
			continue
		}

		if err != redis.ErrNil {
			return vars, err
		}
	}
	return vars, nil
}

// WatchPrefix is not yet implemented.
func (c *Client) WatchPrefix(prefix string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	<-stopChan
	return 0, nil
}
