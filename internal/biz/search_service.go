package biz

import (
	"context"
	"encoding"
	"encoding/base64"
	"encoding/json"
)

type EsSearchUseCase struct {
	repo EsSearchRepo
}

func NewEsSearchUseCase(repo EsSearchRepo) *EsSearchUseCase {
	return &EsSearchUseCase{
		repo: repo,
	}
}

func (e *EsSearchUseCase) SearchProduct(ctx context.Context, query map[string]interface{}, pageToken *PageToken) (*PageResult, error) {
	if len(query) == 0 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	}
	if len(pageToken.NextPageParam) != 0 {
		searchAfter := make([]interface{}, 0)
		err := json.Unmarshal([]byte(pageToken.NextPageParam), &searchAfter)
		if err != nil {
			//s.log.Errorf("search after param unmarshal error: %v", err)
		} else {
			query["search_after"] = searchAfter
		}
	}
	query["size"] = pageToken.Size
	search, err := e.repo.Search(ctx, "product", query)
	return &PageResult{
		Result: search,
		NextToken: &PageToken{
			Size:          10,
			NextPageParam: "",
		},
	}, err
}

type EsSearchRepo interface {
	Search(ctx context.Context, index string, query map[string]interface{}) (*Result, error)
	Update(ctx context.Context, index, id string, content map[string]interface{}) error
	Create(ctx context.Context, index, id string, content map[string]interface{}) error
	Upsert(ctx context.Context, index, id string, content interface{}) error
}

type EsSearchCondition struct {
	Index   string                   // 索引
	Keyword string                   // 关键字
	Query   encoding.BinaryMarshaler // 查询条件
}

type PageToken struct {
	Size          uint32 `json:"size"`
	NextPageParam string `json:"next_page_param"`
}

func (p *PageToken) String() string {
	marshal, err := json.Marshal(p)
	if err != nil {
		//todo: 输出日志
		return ""
	}
	return base64.StdEncoding.EncodeToString(marshal)
}

type PageResult struct {
	Result    *Result
	NextToken *PageToken
}

type Result struct {
	List []interface{}
}
