package data

import (
	"context"
	"fmt"
	"os"
	"testing"
	"yy-shop/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

func Test_itemRepo_FetchByItemName(t *testing.T) {
	logger := log.With(log.NewStdLogger(os.Stdout))
	data, f, err := NewData(&conf.Data{
		Database: &conf.Data_Database{
			Source: "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		},
	}, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer f()
	repo := NewItemRepo(data, logger)
	fetch, err := repo.FetchByItemName(context.Background(), "o", 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range fetch {
		fmt.Printf("%+v\n", item)
	}
}
