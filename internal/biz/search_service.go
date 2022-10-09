package biz

import "context"

type SearchService interface {
	Search(ctx context.Context, condition SearchConditionBuilder) ([]*Result, error)
	GetConditionBuilder() SearchConditionBuilder
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

type Result struct {
	ID *ID
}

type ResultList []*Result

func (r ResultList) GetAllID() []*ID {
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
