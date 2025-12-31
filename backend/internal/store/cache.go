package store

import "github.com/kmdkuk/nfrecap/internal/model"

type Cache interface {
	Get(workTitle string, typ string) (model.Metadata, bool, error)
	Put(workTitle string, typ string, md model.Metadata) error
}
