package conf

import (
	"fmt"
	"github.com/favar/conf/core"
)

type config struct {
	providers []core.ConfigurationProvider
}

func (c *config) GetProviderInfo() []string {
	infos := make([]string, len(c.providers))
	for i, provider := range c.providers {
		infos[i] = provider.String()
	}
	return infos
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

func DefaultConfiguration() core.Configuration {
	return &config{
		providers: make([]core.ConfigurationProvider, 0),
	}
}
