package conf

import (
	"fmt"
	"github.com/favar/conf/core"
	"github.com/favar/conf/provider"
	"github.com/favar/conf/source"
	"io/ioutil"
	"path"
)

type fileProviderAdapter struct {
	adapters map[string]func(source core.ConfigurationSource) (core.ConfigurationProvider, error)
}

func (f fileProviderAdapter) Adapter(source core.ConfigurationSource) (core.ConfigurationProvider, error) {
	ext := path.Ext(source.Id())
	if fn, ok := f.adapters[ext]; ok {
		return fn(source)
	} else {
		return nil, fmt.Errorf("can not found adapter by %s", source.Id())
	}
}

func NewFileProviderAdapter() ProviderAdapter {
	return &fileProviderAdapter{
		adapters: map[string]func(source core.ConfigurationSource) (core.ConfigurationProvider, error){
			".json": func(source core.ConfigurationSource) (core.ConfigurationProvider, error) {
				return provider.JsonProvider(source)
			},
		},
	}
}

func (b *ConfigurationBuilder) AddFiles(paths ...string) *ConfigurationBuilder {
	for _, filePath := range paths {
		fileName := path.Base(filePath)
		if fileName[0:2] == "*." {
			dir := path.Dir(filePath)
			fis, err := ioutil.ReadDir(dir)
			if err != nil {
				panic(err)
			}
			for _, fi := range fis {
				if fi.IsDir() {
					// ignore dir
					continue
				}
				fullPath := path.Join(dir, fi.Name())
				sb := SourceBuilder{
					Source:          source.FileSource(fullPath),
					ProviderAdapter: NewFileProviderAdapter(),
					reloadChanged:   false,
				}
				b.Add(sb)
			}
		} else {
			sb := SourceBuilder{
				Source:          source.FileSource(filePath),
				ProviderAdapter: NewFileProviderAdapter(),
				reloadChanged:   false,
			}
			b.Add(sb)
		}
	}
	return b
}

func (b *ConfigurationBuilder) Include(dir string) *ConfigurationBuilder {

	var recursion func(dir string)
	recursion = func(dir string) {
		fis, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err)
		}
		for _, fi := range fis {
			fullPath := path.Join(dir, fi.Name())
			if fi.IsDir() {
				recursion(fullPath)
			} else {
				b.AddFiles(fullPath)
			}
		}
	}
	recursion(dir)
	return b
}
