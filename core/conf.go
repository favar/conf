package core

type NotFoundError struct {
	desc string
}

func (err NotFoundError) Error() string {
	return err.desc
}

type ConfigurationSource interface {
	Id() string
	SourceId() string
	String() string
}

type ConfigurationProvider interface {
	Load() error
	Get(key string, addr interface{}) error
	String() string
}

type Configuration interface {
	Get(key string, ref interface{}) error
	AddProvider(provider ConfigurationProvider)
	GetProviderInfo() []string
}

type ReloadError interface {
	Catch(err error)
}
