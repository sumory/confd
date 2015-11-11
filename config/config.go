package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/sumory/confd/processor"
	"github.com/sumory/confd/store"
)

var (
	configFile        = ""
	defaultConfigFile = "/data/confd/data/config.toml"
	storeType         string
	connectAddr       string
	confDir           string
	interval          int
	debug             bool
)

//confd
type Config struct {
	Store       string `toml:"store"`
	ConnectAddr string `toml:"connect_addr"`

	ConfDir string `toml:"confdir"` //confd配置文件、模板文件、meta文件目录

	Interval int `toml:"interval"`
	Debug    bool
}

func init() {
	flag.StringVar(&storeType, "store-type", "file", "backend store to use")
	flag.StringVar(&confDir, "confdir", "/data/confd", "confd conf directory")
	flag.StringVar(&configFile, "config-file", "", "the confd config file")
	flag.StringVar(&connectAddr, "connect-addr", "", "backend store address")
	flag.IntVar(&interval, "interval", 600, "backend polling interval")
	flag.BoolVar(&debug, "debug", false, "debug mode")
}

func InitConfig(configDirFromBuild string) (error, *Config, *processor.TemplateConfig, *store.StoreConfig) {
	// 默认配置
	config := Config{
		Store:    "file",
		ConfDir:  "/data/confd",
		Interval: 600,
	}

	if configDirFromBuild != "" { //尝试使用build脚本的变量初始化configFile
		fmt.Printf("configDirFromBuild:%s\n", configDirFromBuild)
		configFileFromBuild := fmt.Sprintf("%s/data/config.toml", configDirFromBuild)
		fmt.Printf("use config file[%s] from build script\n", configFileFromBuild)
		if _, err := os.Stat(configFileFromBuild); !os.IsNotExist(err) {
			configFile = configFileFromBuild
			config.ConfDir = configDirFromBuild
		} else {
			fmt.Printf("use config file[%s] from build script error, file not exist, skip...\n", configFileFromBuild)
		}
	}

	if configFile == "" {
		if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
			fmt.Printf("use defaultConfigFile[%s]\n", defaultConfigFile)
			configFile = defaultConfigFile
		}
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
		Store:       config.Store,
		ConnectAddr: config.ConnectAddr,
	}
	// 模板设置BT
	templateConfig := processor.TemplateConfig{
		ConfDir:     config.ConfDir,
		MetaDir:     filepath.Join(config.ConfDir, "meta"),
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
