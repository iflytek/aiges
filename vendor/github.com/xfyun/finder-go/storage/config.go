package storage

type StorageConfig struct {
	Name   string
	Params map[string]string

	ConfigRootPath  string
	ServiceRootPath string
}
