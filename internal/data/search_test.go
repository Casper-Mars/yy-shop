package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/durationpb"
	"os"
	"testing"
	"time"
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
	param := map[string]interface{}{
		"doc": map[string]interface{}{
			"price": 50,
		},
	}
	bytess, err := json.Marshal(param)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytess))
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:  data.es,
		Index:   "product",
		Refresh: "true",
	})
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(bytess)
	err = indexer.Add(context.Background(), esutil.BulkIndexerItem{
		Action:     "update",
		DocumentID: "1010",
		Body:       reader,
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, item2 esutil.BulkIndexerResponseItem, err error) {
			if err != nil {
				fmt.Println(err)
			}
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = indexer.Close(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
