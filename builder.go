package conf

import (
	"github.com/favar/conf/core"
	"github.com/favar/conf/reload_err"
	"os"
	"sync"
	"time"
)

type ProviderAdapter interface {
	Adapter(source core.ConfigurationSource) (core.ConfigurationProvider, error)
}

type ConfigurationSetting struct {
	Catcher core.ReloadError
}

type ConfigurationBuilder struct {
	setting  *ConfigurationSetting
	conf     core.Configuration
	builders map[string]SourceBuilder
}

type SourceBuilder struct {
	Source          core.ConfigurationSource
	ProviderAdapter ProviderAdapter
	reloadChanged   bool
	lastModTime     time.Time
	provider        core.ConfigurationProvider
}

func (s *SourceBuilder) Build() (provider core.ConfigurationProvider, err error) {

	provider, err = s.ProviderAdapter.Adapter(s.Source)
	s.provider = provider
	return provider, err
}

func (s SourceBuilder) Reload() (err error) {
	if s.reloadChanged {
		var fi os.FileInfo
		fi, err = os.Stat(s.Source.Id())
		if err != nil {
			return
		}
		if s.lastModTime != fi.ModTime() {
			err = s.provider.Load()
		}
	}
	return err
}

var mutex sync.Mutex

func (b *ConfigurationBuilder) Add(builder SourceBuilder) {
	source := builder.Source
	if source == nil {
		panic("Source is nil")
	}
	id := source.Id()

	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := b.builders[id]; !ok {
		b.builders[id] = builder
	}
}

func (b *ConfigurationBuilder) monitor() {
	go func() {
		for true {
			// 500ms check once
			time.Sleep(500 * time.Millisecond)
			for _, builder := range b.builders {
				err := builder.Reload()
				if err != nil && b.setting.Catcher != nil {
					b.setting.Catcher.Catch(err)
				}
			}
		}
	}()
}

func (b *ConfigurationBuilder) Build() (core.Configuration, error) {

	for _, builder := range b.builders {
		if provider, err := builder.Build(); err == nil {
			b.conf.AddProvider(provider)
		} else {
			return nil, err
		}
	}
	b.monitor()
	return b.conf, nil
}

func (b *ConfigurationBuilder) OnSetting(x func(setting *ConfigurationSetting)) *ConfigurationBuilder {
	x(b.setting)
	return b
}

func Builder() *ConfigurationBuilder {
	return &ConfigurationBuilder{
		builders: make(map[string]SourceBuilder),
		conf:     DefaultConfiguration(),
		setting: &ConfigurationSetting{
			Catcher: reload_err.FmtReloadError(),
		},
	}
}
