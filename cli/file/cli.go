package file

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
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

func NewFileCli(filePath string) (*Cli, error) {
	var fileStore *FileStore
	_, err := toml.DecodeFile(filePath, &fileStore)
	fmt.Printf("%s %v\n", filePath, fileStore)

	if err != nil {
		return nil, fmt.Errorf("Cannot process file store %s - %s", filePath, err.Error())
	}

	return &Cli{
		FileStore: fileStore,
	}, nil
}

func (c *Cli) GetValues(keys []string) (map[string]interface{}, error) {
	return c.FileStore.Data, nil
}

func (c *Cli) SetValues(values map[string]interface{}) error {

	filename := c.FileStore.Path

	// 创建临时文件
	temp, err := ioutil.TempFile(filepath.Dir(filename), "."+filepath.Base(filename))
	defer os.Remove(temp.Name())
	defer temp.Close()

	w := bufio.NewWriter(temp)
	toml.NewEncoder(w).Encode(values)
	os.Chmod(temp.Name(), 0644)

	err = os.Rename(temp.Name(), filename)
	if err != nil {
		return err
	}

	return nil
}
