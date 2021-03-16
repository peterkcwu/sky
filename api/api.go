package api

import (
	"github.com/new_web/config"
	"github.com/new_web/store"
	"gopkg.in/resty.v1"
)

type ApiClient struct {
	Conf           *config.Config
	Store          *store.Store
	RestyClient    *resty.Client
}

func NewApiClient(conf *config.Config) (*ApiClient, error) {
	//newStore, err := store.NewStore(conf)
	//if err != nil {
	//	return nil, err
	//}
	clientR := resty.New()
	clientR.SetDebug(true)
	return &ApiClient{Conf: conf, RestyClient: clientR}, nil

}