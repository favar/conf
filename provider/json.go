package provider

import (
	"encoding/json"
	"fmt"
	"github.com/favar/conf/core"
	"os"
)

type provider struct {
	source core.ConfigurationSource
	parser core.Parser
}

func (p *provider) String() string {
	return p.source.String()
}

func (p *provider) Load() error {
	file, err := os.Open(p.source.Id())
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var v interface{}
	err = decoder.Decode(&v)
	if s, ok := v.(map[string]interface{}); ok {
		p.parser = core.JsonParser(s)
	} else {
		err = fmt.Errorf("path[%s] is not map[string]interface{}", p.source.Id())
	}
	return err
}

func (p *provider) Get(key string, addr interface{}) error {
	return p.parser.Get(key, addr)
}

func JsonProvider(source core.ConfigurationSource) (pro core.ConfigurationProvider, err error) {
	pro = &provider{source: source}
	err = pro.Load()
	return pro, err
}
