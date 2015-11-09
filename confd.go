package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/sumory/confd/config"
	"github.com/sumory/confd/store"
	"github.com/sumory/confd/processor"
	"os"
	"os/signal"
	"syscall"
)


func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     false,
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05.00000",
	})
	log.SetLevel(log.DebugLevel)
}

var configDirFromBuild string = ""

func main() {
	flag.Parse()

	err, myConfig, templateConfig, storeConfig := config.InitConfig(configDirFromBuild)
	if err != nil {
		log.Error(err.Error())
	}

	log.Info("Starting confd")

	storeClient, err := store.New(storeConfig)
	if err != nil {
		log.Error(err.Error())
	}

	templateConfig.StoreClient = storeClient

	if myConfig.Debug {
		processor.Process(templateConfig)
		return
	}

	stopChan := make(chan bool)
	doneChan := make(chan bool)
	errChan := make(chan error, 10)

	processor := processor.NewIntervalProcessor(templateConfig, stopChan, doneChan, errChan, myConfig.Interval)

	go processor.Process()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case err := <-errChan:
			log.Error(err.Error())
		case s := <-signalChan:
			log.Info(fmt.Sprintf("Captured %v. Exiting...", s))
			close(doneChan)
		case <-doneChan:
			os.Exit(0)
		}
	}
}
