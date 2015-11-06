package processor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/sumory/confd/store"
	"github.com/sumory/confd/utils"
)

type MetaConfig struct {
	MetaObject MetaObject `toml:"metaObject"`
}

type MetaObject struct {
	Tmpl        string //模板文件，用于生成最后的配置文件
	Dest        string //最终配置文件的路径
	Keys        []string //需要使用的key

	FileMode    os.FileMode
	Mode        string
	funcMap     map[string]interface{}
	storeClient store.StoreClient

	kvStore utils.KVStore
}

var EmptyErr = errors.New("empty meta file")

func NewMetaObject(path string, config *TemplateConfig) (*MetaObject, error) {
	if config.StoreClient == nil {
		return nil, errors.New("StoreClient is not found.")
	}

	log.Debug("Loading meta info from ", path)

	var tc *MetaConfig
	_, err := toml.DecodeFile(path, &tc)
	if err != nil {
		return nil, fmt.Errorf("Cannot process meta file %s - %s", path, err.Error())
	}

	log.Debug(fmt.Sprintf("metaObject: %+v", tc.MetaObject))
	tr := tc.MetaObject
	tr.storeClient = config.StoreClient
	tr.kvStore = utils.NewKVStore()
	tr.funcMap = tr.kvStore.FuncMap

	if tr.Tmpl == "" {
		return nil, EmptyErr
	}
	tr.Tmpl = filepath.Join(config.TemplateDir, tr.Tmpl)

	return &tr, nil
}

func (t *MetaObject) process() error {
	if err := t.setFileMode(); err != nil {
		return err
	}

	result, err := t.setVars()
	if err != nil {
		return err
	}
	if err = t.createConfigFile(result); err != nil {
		return err
	}

	return nil
}

// setVars sets the Vars for template resource.
func (t *MetaObject) setVars() (map[string]interface{}, error) {
	var err error
	log.Debug("Retrieving keys from store")
	result, err := t.storeClient.GetValues(t.Keys)


	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", result)
	t.kvStore.Clean()
	for k, v := range result {
		t.kvStore.Set(k, v)
	}

	return result, nil
}

func (t *MetaObject) createConfigFile(data map[string]interface{}) error {
	log.Debug("Using template " + t.Tmpl)

	if !utils.IsFileExist(t.Tmpl) {
		return errors.New("Missing meta file: " + t.Tmpl)
	}

	log.Debug("Compiling source template " + t.Tmpl)
	tmpl, err := template.New(path.Base(t.Tmpl)).Funcs(t.funcMap).ParseFiles(t.Tmpl)
	if err != nil {
		return fmt.Errorf("Unable to process template file %s, %s", t.Tmpl, err)
	}

	// 创建临时文件
	temp, err := ioutil.TempFile(filepath.Dir(t.Dest), "."+filepath.Base(t.Dest))
	defer os.Remove(temp.Name())
	defer temp.Close()

	log.Debug("Temp path: ", temp.Name())

	if err != nil {
		return err
	}

	if err = tmpl.Execute(temp, nil); err != nil {
		temp.Close()
		os.Remove(temp.Name())
		return err
	}
	os.Chmod(temp.Name(), t.FileMode)

	log.Debug("Overwriting target config: ", t.Dest)
	err = os.Rename(temp.Name(), t.Dest)
	if err != nil {
		log.Fatal("Rename ", temp.Name(), " to ", t.Dest, " failed")
		return err
	}

	log.Info("Target config has been updated: ", t.Dest)
	return nil
}

func (t *MetaObject) setFileMode() error {
	if t.Mode == "" {
		if !utils.IsFileExist(t.Dest) {
			t.FileMode = 0644
		} else {
			fi, err := os.Stat(t.Dest)
			if err != nil {
				return err
			}
			t.FileMode = fi.Mode()
		}
	} else {
		mode, err := strconv.ParseUint(t.Mode, 0, 32)
		if err != nil {
			return err
		}
		t.FileMode = os.FileMode(mode)
	}
	return nil
}
