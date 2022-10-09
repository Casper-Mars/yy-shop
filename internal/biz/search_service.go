package biz

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"
)

type SearchMgr interface {
	Search(ctx context.Context, condition SearchConditionBuilder) (ResultList, error)
	GetConditionBuilder() SearchConditionBuilder
	SearchByPage(ctx context.Context, condition SearchConditionBuilder, pageToken PageToken) (*PageResult, error)
}

type PageToken struct {
	size           uint32
	sortUpdateTime time.Time
	sortID         uint32
}

func NewPageToken(size uint32, sortUpdateTime time.Time, sortID uint32) *PageToken {
	return &PageToken{
		size:           size,
		sortUpdateTime: sortUpdateTime,
		sortID:         sortID,
	}
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
	result    ResultList
	pageToken *PageToken
}

func (p *PageResult) GetResult() ResultList {
	return p.result
}

func (p *PageResult) PageToke() *PageToken {
	return p.pageToken
}

type ID struct {
	_int64  int64
	_uint64 uint64
	_int32  int32
	_uint32 uint32
	_string string
}

func (i ID) Int64() int64 {
	return i._int64
}

func (i ID) Uint64() uint64 {
	return i._uint64
}

func (i ID) Int32() int32 {
	return i._int32
}

func (i ID) Uint32() uint32 {
	return i._uint32
}

func (i ID) String() string {
	return i._string
}

type IDList []*ID

func (i IDList) AsInt64() []int64 {
	if len(i) == 0 {
		return []int64{}
	}
	ids := make([]int64, len(i))
	for i, id := range i {
		ids[i] = id.Int64()
	}
	return ids
}

func (i IDList) AsUint64() []uint64 {
	if len(i) == 0 {
		return []uint64{}
	}
	ids := make([]uint64, len(i))
	for i, id := range i {
		ids[i] = id.Uint64()
	}
	return ids
}

func (i IDList) AsInt32() []int32 {
	if len(i) == 0 {
		return []int32{}
	}
	ids := make([]int32, len(i))
	for i, id := range i {
		ids[i] = id.Int32()
	}
	return ids
}

func (i IDList) AsUint32() []uint32 {
	if len(i) == 0 {
		return []uint32{}
	}
	ids := make([]uint32, len(i))
	for i, id := range i {
		ids[i] = id.Uint32()
	}
	return ids
}

func (i IDList) AsString() []string {
	if len(i) == 0 {
		return []string{}
	}
	ids := make([]string, len(i))
	for i, id := range i {
		ids[i] = id.String()
	}
	return ids
}

type Result struct {
	ID *ID
}

type ResultList []*Result

func (r ResultList) GetAllID() IDList {
	if len(r) == 0 {
		return []*ID{}
	}
	ids := make([]*ID, len(r))
	for i, result := range r {
		ids[i] = result.ID
	}
	return ids
}

type SearchConditionBuilder interface {
	SetKeyWord(keyWord string, value string) SearchConditionBuilder
	Build() *interface{}
}
