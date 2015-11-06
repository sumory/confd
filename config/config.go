package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/sumory/confd/store"
	"github.com/sumory/confd/processor"
	"fmt"
)

var (
	configFile = ""
	defaultConfigFile = "/data/confd/data/config.toml"
	storeType string
	connectAddr string
	confDir string
	interval int
	debug bool
)

//confd
type Config struct {
	Store       string `toml:"store"`
	ConnectAddr string `toml:"connect_addr"`

	ConfDir     string `toml:"confdir"` //confd模板文件、meta文件目录

	Interval    int    `toml:"interval"`
	Debug       bool
}





func init() {
	flag.StringVar(&storeType, "store-type", "redis", "backend to use")
	flag.StringVar(&confDir, "confdir", "/data/confd", "confd conf directory")
	flag.StringVar(&configFile, "config-file", "", "the confd config file")
	flag.IntVar(&interval, "interval", 600, "backend polling interval")
	flag.StringVar(&connectAddr, "connect-addr", "", "list of backend nodes")
	flag.BoolVar(&debug, "debug", false, "debug mode")
}

func InitConfig() (error, *Config, *processor.TemplateConfig, *store.StoreConfig) {
	if configFile == "" {
		if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
			configFile = defaultConfigFile
		}
	}

	// 默认配置
	config := Config{
		Store:    "redis",
		ConfDir:  "/data/confd",
		Interval: 600,
	}

	//从toml文件更新配置
	if configFile == "" {
		log.Debug("Skipping confd config file.")
	} else {
		log.Debug("Loading " + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err, nil, nil, nil
		}
		_, err = toml.Decode(string(configBytes), &config)
		if err != nil {
			return err, nil, nil, nil
		}
	}

	// 根据命令行参数更新配置
	processFlags(&config)

	if config.ConnectAddr == "" {
		switch config.Store {
		case "file":
			config.ConnectAddr = fmt.Sprintf("%s/data/filestore.toml", config.ConfDir)
		case "redis":
			config.ConnectAddr = "127.0.0.1:6379"
		case "zookeeper":
			config.ConnectAddr = "127.0.0.1:2181"
		}
	}

	log.Info("Store set to " + config.Store)

	//后端client设置
	storeConfig := store.StoreConfig{
		Store:config.Store,
		ConnectAddr:config.ConnectAddr,
	}
	// 模板设置BT
	templateConfig := processor.TemplateConfig{
		ConfDir:     config.ConfDir,
		MetaDir:   filepath.Join(config.ConfDir, "meta"),
		TemplateDir: filepath.Join(config.ConfDir, "templates"),
	}

	return nil, &config, &templateConfig, &storeConfig
}

func processFlags(config *Config) {
	flag.Visit(func(f *flag.Flag) {
		setConfigFromFlag(f, config)
	})
}

func setConfigFromFlag(f *flag.Flag, config *Config) {
	switch f.Name {
	case "store-type":
		config.Store = storeType
	case "confdir":
		config.ConfDir = confDir
	case "connect-addr":
		config.ConnectAddr = connectAddr
	case "interval":
		config.Interval = interval
	case "debug":
		config.Debug = debug
	}
}
