package data

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"yy-shop/internal/biz"
)

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

func (s *searchRepo) Search(ctx context.Context, condition biz.EsSearchCondition) (biz.ResultList, error) {
	//TODO implement me
	panic("implement me")
}

func (s *searchRepo) SearchByPage(ctx context.Context, condition *biz.EsSearchCondition, pageToken *biz.PageToken) (*biz.PageResult, error) {
	query := map[string]interface{}{
		"query": condition.Query,
		"size":  pageToken.Size,
		"sort": map[string]string{
			//"update_time": "desc",
			"id": "desc",
		},
	}
	if len(pageToken.NextPageParam) != 0 {
		query["search_after"] = pageToken.NextPageParam
	}
	s.log.Infof("query: %v", query)
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		s.log.Errorf("json encode error: %v", err)
		return nil, err
	}
	response, err := s.data.es.Search(
		s.data.es.Search.WithContext(ctx),
		s.data.es.Search.WithIndex("product"),
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
	return nil, nil
}
