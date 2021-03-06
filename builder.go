package conf

import (
	"fmt"
	"github.com/favar/conf/core"
	"os"
	"sync"
	"time"
)

type ConfigurationBuilder struct {
	conf     core.Configuration
	builders map[string]SourceBuilder
}

type SourceBuilder struct {
	Source        core.ConfigurationSource
	Build         func(source core.ConfigurationSource) (core.ConfigurationProvider, error)
	ReloadChanged bool
	lastModTime   time.Time
}

var mutex sync.Mutex

func (b *ConfigurationBuilder) Add(builder SourceBuilder) {
	source := builder.Source
	if source == nil {
		panic("Source is nil")
	}
	id := source.SourceId()

	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := b.builders[id]; !ok {
		b.builders[id] = builder
	}
}

func (b ConfigurationBuilder) monitor() {
	monitorBuilders := make([]SourceBuilder, 0)
	for _, builder := range b.builders {
		if builder.ReloadChanged {
			monitorBuilders = append(monitorBuilders, builder)
		}
	}
	if len(monitorBuilders) == 0 {
		return
	}
	go func() {
		for true {
			// 500ms check once
			time.Sleep(500 * time.Millisecond)
			for _, builder := range monitorBuilders {
				if fi, err := os.Stat(builder.Source.Path()); err == nil {
					last := fi.ModTime()
					if last != builder.lastModTime {
						builder.lastModTime = last
						if provider, err := builder.Build(builder.Source); err == nil {
							b.conf.AddProvider(provider)
						} else {
							// todo
							//return nil, err
						}
					}
				}
			}
		}
	}()
}

func (b *ConfigurationBuilder) Build() (core.Configuration, error) {

	for _, builder := range b.builders {
		if provider, err := builder.Build(builder.Source); err == nil {
			b.conf.AddProvider(provider)
			if fi, err := os.Stat(builder.Source.Path()); err == nil {
				builder.lastModTime = fi.ModTime()
			}
		} else {
			return nil, err
		}
	}
	// todo
	//b.monitor()
	return b.conf, nil
}

func Builder() *ConfigurationBuilder {
	conf := &config{
		providers: make([]core.ConfigurationProvider, 0),
	}
	return &ConfigurationBuilder{builders: make(map[string]SourceBuilder), conf: conf}
}

type config struct {
	providers []core.ConfigurationProvider
}

func (c *config) AddProvider(provider core.ConfigurationProvider) {
	c.providers = append(c.providers, provider)
}

func (c *config) Get(key string, ref interface{}) error {
	for _, provider := range c.providers {
		if err := provider.Get(key, ref); err == nil {
			break
		} else {
			// if not found, find next
			if _, ok := err.(core.NotFoundError); ok {
				fmt.Println(err)
				continue
			} else {
				return err
			}
		}
	}
	return nil
}
