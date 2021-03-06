package conf

import (
	"github.com/favar/conf/core"
	"github.com/favar/conf/provider"
	"github.com/favar/conf/source"
)

//func (s jsonSource) Build() (ConfigurationProvider, error) {
//	p := provider{path: s.path}
//	err := p.Load()
//	if s.ReloadOnChanged() {
//		info, _ := os.Stat(s.path)
//		s.lastModTime = info.ModTime()
//
//		go func() {
//			for true {
//				time.Sleep(time.Duration(s.reloadDelay * 1000000000))
//				if info2, errf := os.Stat(s.path); errf == nil {
//					// todo err log
//					last := info2.ModTime()
//					if last != s.lastModTime {
//						s.lastModTime = last
//						p.Load()
//					}
//				}
//			}
//		}()
//	}
//
//	return &p, err
//}

func (b *ConfigurationBuilder) AddJsonFiles(paths ...string) *ConfigurationBuilder {
	for _, path := range paths {
		sb := SourceBuilder{
			Source: source.FileSource(path),
			Build: func(source core.ConfigurationSource) (core.ConfigurationProvider, error) {
				return provider.JsonProvider(source)
			},
			ReloadChanged: false,
		}
		b.Add(sb)
	}
	return b
}

func (b *ConfigurationBuilder) AddJsonFolder(paths ...string) *ConfigurationBuilder {
	return b
}
