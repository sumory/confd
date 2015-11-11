package file

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileStore struct {
	Path string
	Data map[string]interface{} `toml:"data"`
}

type Cli struct {
	FileStore *FileStore
}

var EmptyErr = errors.New("empty key/value")

func NewFileCli(filePath string) (error, *Cli) {
	var fileStore *FileStore
	_, err := toml.DecodeFile(filePath, &fileStore)
	log.Infof("%s %v", filePath, fileStore)

	if err != nil {
		return fmt.Errorf("Cannot process file store %s - %s", filePath, err.Error()), nil
	}

	fileStore.Path = filePath
	return nil, &Cli{
		FileStore: fileStore,
	}
}

//重新获取所有value
func (c *Cli) Fetch() error {
	_, err := toml.DecodeFile(c.FileStore.Path, &c.FileStore)
	//log.Infof("Fetch %s, values: %+v", c.FileStore.Path, *c.FileStore)
	if err != nil {
		return fmt.Errorf("Cannot process file store %s - %s", c.FileStore.Path, err.Error())
	}

	return nil
}

func (c *Cli) GetAll() (map[string]interface{}, error) {
	return c.FileStore.Data, nil
}

func (c *Cli) GetValues(keys []string) (map[string]interface{}, error) {
	return c.FileStore.Data, nil
}

func (c *Cli) GetValue(key string) (interface{}, error) {
	if value, ok := c.FileStore.Data[key]; ok {
		return value, nil
	}
	return nil, EmptyErr
}

func (c *Cli) DeleteKey(key string) error {
	err := c.Fetch()
	if err != nil {
		return err
	}

	delete(c.FileStore.Data, key)
	return c.setValues(c.FileStore.Data)
}

func (c *Cli) setValues(values map[string]interface{}) error {

	filename := c.FileStore.Path

	// 创建临时文件
	temp, err := ioutil.TempFile(filepath.Dir(filename), "."+filepath.Base(filename))
	temp.WriteString("[data]\n")
	defer os.Remove(temp.Name())
	defer temp.Close()

	//	fmt.Println(filename)
	//	fmt.Println(temp.Name())
	//	fmt.Println(values)

	w := bufio.NewWriter(temp)
	toml.NewEncoder(w).Encode(values)
	os.Chmod(temp.Name(), 0644)

	err = os.Rename(temp.Name(), filename)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cli) SetValue(key string, value interface{}) error {
	err := c.Fetch()
	if err != nil {
		return err
	}

	c.FileStore.Data[key] = value
	return c.setValues(c.FileStore.Data)
}
