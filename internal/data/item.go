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
	err = i.table.WithContext(ctx).Where("item_name like ?", fmt.Sprintf("%%%s%%", itemName)).Limit(int(limit)).Offset(int(pageToken)).Find(&out).Error
	if err != nil {
		i.log.Errorf("FetchByItemName failed to Find, err:%v", err)
		return out, err
	}
	return out, nil
}
