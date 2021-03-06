package source

import (
	"github.com/favar/conf/core"
	"strings"
)

type source struct {
	path string
}

func (s source) String() string {
	return s.path
}

func (s source) SourceId() string {
	return strings.ToLower(s.path)
}

func (s source) Id() string {
	return s.path
}

func FileSource(path string) core.ConfigurationSource {
	return source{path}
}
