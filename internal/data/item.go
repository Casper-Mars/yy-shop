package data

import (
	"context"
	"fmt"
	"yy-shop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type itemRepo struct {
	data  *Data
	log   *log.Helper
	table *gorm.DB
}

func NewItemRepo(data *Data, logger log.Logger) biz.ItemRepo {
	return &itemRepo{
		data:  data,
		log:   log.NewHelper(logger),
		table: data.db.Table("item"),
	}
}

func (i *itemRepo) FetchByItemName(ctx context.Context, itemName string, pageToken, pageSize uint32) (itemInfoList []*biz.ItemInfo, err error) {
	limit := uint32(30)
	if pageSize < 30 {
		limit = pageSize
	}
	out := make([]*biz.ItemInfo, 0)
	err = i.table.WithContext(ctx).Where("id > ? and item_name like ?", pageToken, fmt.Sprintf("%%%s%%", itemName)).Limit(int(limit)).Find(&out).Error
	if err != nil {
		i.log.Errorf("FetchByItemName failed to Find, err:%v", err)
		return out, err
	}
	return out, nil
}

func (i *itemRepo) FetchByIds(ctx context.Context, ids ...uint32) (itemInfoList biz.ItemList, err error) {
	if len(ids) == 0 {
		return nil, biz.ErrNoResult
	}
	err = i.table.WithContext(ctx).Where("id in (?)", ids).Find(&itemInfoList).Error
	if err != nil {
		i.log.Errorf("FetchByIds failed to Find, err:%v", err)
		return nil, err
	}
	return
}
