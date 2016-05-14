package myproxy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type ProxyConfig struct {
	ListenHost string
	BackHost   string
}

func NewConfig(fpath string) *ProxyConfig {
	log.Println("will load config from file:", fpath)

	res := ProxyConfig{}
	fin, err := os.Open(fpath)
	if err != nil {
		log.Println("load config from file error:", err)
		return nil
	}
	defer fin.Close()

	buff, err := ioutil.ReadAll(fin)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(buff, &res)
	if err != nil {
		return nil
	}

	log.Printf("config loaded: %#v\n", res)
	return &res
}
