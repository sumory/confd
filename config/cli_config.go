package config

import (
	"io/ioutil"
	"os"

	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/sumory/confd/cli"
)

func InitCliConfig(cliStoreType, cliConfDir, cliConfigFile, clieConnectAddr string, cliConfig *cli.CliConfig) error {
	if cliConfigFile == "" {
		if _, err := os.Stat(defaultConfigFile); !os.IsNotExist(err) {
			cliConfigFile = defaultConfigFile
		}
	}

	// 默认配置
	config := Config{
		Store:   "file",
		ConfDir: "/data/confd",
	}

	//从toml文件更新配置
	if cliConfigFile == "" {
		log.Debug("Skipping confd config file.")
	} else {
		log.Debug("Loading " + cliConfigFile)
		configBytes, err := ioutil.ReadFile(cliConfigFile)
		if err != nil {
			return err
		}
		_, err = toml.Decode(string(configBytes), &config)
		if err != nil {
			return err
		}
	}

	//从命令行参数覆盖配置
	config.Store = cliStoreType
	config.ConfDir = cliConfDir
	config.ConnectAddr = clieConnectAddr

	//配置connectAddr
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

	//cli设置
	cliConfig.Store = config.Store
	cliConfig.ConnectAddr = config.ConnectAddr

	return nil
}
