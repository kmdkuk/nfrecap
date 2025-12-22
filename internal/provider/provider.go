package provider

import "github.com/kmdkuk/nfrecap/internal/model"

type Provider interface {
	Lookup(workTitle string, typ string) (model.Metadata, bool, error)
}
