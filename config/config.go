package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

var (
	// ErrNoFileFound returned when there is no config file found
	ErrNoFileFound = errors.New("no config file found")
)

func Read(dest interface{}, paths ...string) error {
	for _, path := range paths {

		// check if this path is exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		// load config
		ext := filepath.Ext(path)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		switch {
		case ext == ".yaml" || ext == ".yml":
			return yaml.Unmarshal(content, dest)
		}
	}
	return ErrNoFileFound
}
