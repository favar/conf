package core

type NotFoundError struct {
	desc string
}

func (err NotFoundError) Error() string {
	return err.desc
}

type ConfigurationSource interface {
	Path() string
	SourceId() string
	//Build() (ConfigurationProvider, error)
}

type ConfigurationProvider interface {
	Load() error
	Get(key string, addr interface{}) error
}

type Configuration interface {
	Get(key string, ref interface{}) error
	AddProvider(provider ConfigurationProvider)
}
