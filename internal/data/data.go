package data

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
	"yy-shop/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewItemRepo)

// Data .
type Data struct {
	db *gorm.DB
	es *elasticsearch.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	db, err := gorm.Open(mysql.Open(c.GetDatabase().GetSource()), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	config := elasticsearch.Config{
		Addresses: c.Elasticsearch.GetAddr(),
		Transport: &http.Transport{
			ResponseHeaderTimeout: c.Elasticsearch.GetTimeout().AsDuration(),
		},
	}
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		panic(err)
	}
	info, err := client.Info()
	if err != nil {
		panic(err)
	}
	fmt.Println(info)

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db: db,
		es: client,
	}, cleanup, nil
}
