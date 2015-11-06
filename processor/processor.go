package processor

import (
	"github.com/sumory/confd/utils"
	"time"
	log "github.com/Sirupsen/logrus"
)

type Processor interface {
	Process()
}

type intervalProcessor struct {
	config   *TemplateConfig
	stopChan chan bool
	doneChan chan bool
	errChan  chan error
	interval int
}

func Process(config *TemplateConfig) error {
	ts, err := getMetaObjects(config)
	if err != nil {
		return err
	}
	return process(ts)
}

func process(ts []*MetaObject) error {
	var lastErr error
	for _, t := range ts {
		if err := t.process(); err != nil {
			log.Error(err.Error())
			lastErr = err
		}
	}
	return lastErr
}

func NewIntervalProcessor(config *TemplateConfig, stopChan, doneChan chan bool, errChan chan error, interval int) Processor {
	return &intervalProcessor{config, stopChan, doneChan, errChan, interval}
}

func (p *intervalProcessor) Process() {
	defer close(p.doneChan)
	for {
		ts, err := getMetaObjects(p.config)
		if err != nil {
			log.Fatal(err.Error())
			break
		}
		process(ts)

		select {
		case <-p.stopChan:
			break
		case <-time.After(time.Duration(p.interval) * time.Second):
			log.Infof("Restart process all configurations.")
			continue
		}
	}
}

func getMetaObjects(config *TemplateConfig) ([]*MetaObject, error) {
	var lastError error
	metaObjects := make([]*MetaObject, 0)

	paths, err := utils.RecursiveFindFiles(config.MetaDir, "*toml")
	if err != nil {
		return nil, err
	}
	for _, p := range paths {
		log.Infof("New meta info from path %s", p)
		t, err := NewMetaObject(p, config)
		if err != nil {
			lastError = err
			continue
		}
		metaObjects = append(metaObjects, t)
	}
	return metaObjects, lastError
}

