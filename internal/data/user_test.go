package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"testing"
	"yy-shop/internal/biz"
	"yy-shop/internal/conf"
)

func TestUserRepo_Save(t *testing.T) {
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
	repo := NewUserRepo(data, logger)
	_, err = repo.Save(context.Background(), &biz.User{
		Username: "123",
		Password: "123",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepo_Fetch(t *testing.T) {
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
	repo := NewUserRepo(data, logger)
	fetch, err := repo.Fetch(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", fetch)
}
