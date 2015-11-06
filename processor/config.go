package processor

import "github.com/sumory/confd/store"

type TemplateConfig struct {
	ConfDir     string
	MetaDir     string
	TemplateDir string
	StoreClient store.StoreClient
}
