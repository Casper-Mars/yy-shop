package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/durationpb"
	"os"
	"testing"
	"time"
	"yy-shop/internal/biz"
	"yy-shop/internal/conf"
)

func Test_searchRepo_SearchByPage(t *testing.T) {
	logger := log.With(log.NewStdLogger(os.Stdout))
	data, f, err := NewData(&conf.Data{
		Database: &conf.Data_Database{
			Source: "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		},
		Elasticsearch: &conf.Data_Elasticsearch{
			Addr:    []string{"http://localhost:9200"},
			Timeout: durationpb.New(2 * time.Second),
		},
	}, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer f()
	esSearchRepo := NewSearchRepo(logger, data)
	page, err := esSearchRepo.SearchByPage(context.Background(), &biz.EsSearchCondition{
		Query: map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		Index: "product",
	}, &biz.PageToken{
		Size: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(page)
}
