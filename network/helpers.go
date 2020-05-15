package network

import (
	"io/ioutil"
	"path/filepath"
)

func ConfigFromFile(name string) (*Config, error) {
	var cfg Config
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return &cfg, cfg.Parse(data)
}

// List networks in directory
func List(directory string) ([]Network, error) {
	abs, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	list, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var ans []Network
	for _, item := range list {
		if item.IsDir() && IsValidName(item.Name()) {
			ans = append(ans, Network{Root: filepath.Join(abs, item.Name())})
		}
	}
	return ans, nil
}
