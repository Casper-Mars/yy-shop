package data

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"yy-shop/internal/biz"
)

type esUpsertBody struct {
	Doc    interface{} `json:"doc"`
	Upsert interface{} `json:"upsert"`
}

type searchRepo struct {
	data *Data
	log  *log.Helper
}

func NewSearchRepo(logger log.Logger, data *Data) biz.EsSearchRepo {
	return &searchRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (s *searchRepo) Search(ctx context.Context, index string, query map[string]interface{}) (*biz.Result, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		s.log.Errorf("json encode error: %v", err)
		return nil, err
	}
	response, err := s.data.es.Search(
		s.data.es.Search.WithContext(ctx),
		s.data.es.Search.WithIndex(index),
		s.data.es.Search.WithBody(&buf),
		s.data.es.Search.WithPretty(),
	)
	if err != nil {
		s.log.Errorf("search error: %v", err)
		return nil, err
	}
	defer response.Body.Close()
	var resp = map[string]interface{}{}
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		s.log.Errorf("decode response error: %v", err)
		return nil, err
	}
	s.log.Infof("response: %v", resp)
	hits, ok := resp["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		s.log.Errorf("convert response error: %v", err)
		return nil, err
	}
	result := &biz.Result{
		List: make([]interface{}, 0, len(hits)),
	}
	for _, hit := range hits {
		result.List = append(result.List, hit)
	}
	return result, nil
}

func (s *searchRepo) Update(ctx context.Context, index, id string, content map[string]interface{}) error {
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"doc": content,
	}
	err := json.NewEncoder(&buf).Encode(doc)
	if err != nil {
		s.log.Errorf("json encode error: %v", err)
		return err
	}
	resp, err := s.data.es.Update(
		index,
		id,
		&buf,
		s.data.es.Update.WithRefresh("true"),
	)
	if err != nil {
		s.log.Errorf("create error: %v", err)
		return err
	}
	defer resp.Body.Close()
	s.log.Debug(resp.String())
	return nil
}

func (s *searchRepo) Create(ctx context.Context, index, id string, content map[string]interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(content)
	if err != nil {
		s.log.Errorf("json encode error: %v", err)
		return err
	}
	resp, err := s.data.es.Create(
		index,
		id,
		&buf,
		s.data.es.Create.WithRefresh("true"),
	)
	if err != nil {
		s.log.Errorf("create error: %v", err)
		return err
	}
	defer resp.Body.Close()
	s.log.Debug(resp.StatusCode)
	s.log.Debug(resp.String())
	return nil
}

func (s *searchRepo) Upsert(ctx context.Context, index, id string, content interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(esUpsertBody{
		Doc:    content,
		Upsert: content,
	})
	if err != nil {
		s.log.Errorf("json encode error: %v", err)
		return err
	}
	resp, err := s.data.es.Update(
		index,
		id,
		&buf,
		s.data.es.Update.WithRefresh("true"),
	)
	if err != nil {
		s.log.Errorf("create error: %v", err)
		return err
	}
	defer resp.Body.Close()
	s.log.Debug(resp.StatusCode)
	s.log.Debug(resp.String())
	return nil
}
