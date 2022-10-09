package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrItemNotExist = errors.New("商品名称有误")
	ErrNoResult     = errors.New("没有搜索结果")
)

type ItemInfoWithSeller struct {
	SellerID       uint32  // 卖家ID
	SellerAvatar   string  // 卖家头像
	SellerNickName string  // 卖家昵称
	ItemId         uint32  // 商品ID
	ItemName       string  // 商品名称
	Price          float64 // 价格
	IconUrl        string  // 商品图片连接
	BookedCnt      uint32  // 想要的人数
}

// ItemInfo 商品信息
type ItemInfo struct {
	ItemId    uint32  `gorm:"column:id"`         // 商品ID
	ItemName  string  `gorm:"column:item_name"`  // 商品名称
	IconUrl   string  `gorm:"column:icon_url"`   // 商品图片连接
	Price     float64 `gorm:"column:price"`      // 价格
	SellerId  uint32  `gorm:"column:seller_id"`  // 卖家id
	BookedCnt uint32  `gorm:"column:booked_cnt"` // 想要的人数
}

type ItemList []*ItemInfo

func (i ItemList) GetSellerIDs() []int64 {
	if len(i) == 0 {
		return []int64{}
	}
	ids := make([]int64, len(i))
	for i, item := range i {
		ids[i] = int64(item.SellerId)
	}
	return ids
}

type ItemRepo interface {
	// FetchByItemName 获取指定用户名的用户的信息，如果用户不存在，则返回 ErrUserNotExist。
	FetchByItemName(ctx context.Context, itemName string, pageToken, pageSize uint32) (itemInfoList []*ItemInfo, err error)
	// FetchByIds 批量获取指定id的商品信息
	FetchByIds(ctx context.Context, ids ...uint32) (itemInfoList ItemList, err error)
}

type ProductCache interface {
}

type ProductMgr struct {
	userRepo UserRepo
	itemRepo ItemRepo
	logger   *log.Helper
}

//NewProductMgr 创建一个AccountUseCase，依赖作为参数传入
func NewProductMgr(logger log.Logger, userRepo UserRepo, producctRepo ItemRepo) *ProductMgr {
	return &ProductMgr{
		userRepo: userRepo,
		itemRepo: producctRepo,
		logger:   log.NewHelper(logger),
	}
}

//SearchItem 搜索商品
func (p *ProductMgr) SearchItem(ctx context.Context, ids ...uint32) ([]*ItemInfoWithSeller, error) {

	out := make([]*ItemInfoWithSeller, 0)
	itemInfoList, err := p.itemRepo.FetchByIds(ctx, ids...)
	if err != nil {
		p.logger.Errorf("SearchItem failed to FetchByItemName, err:%v", err)
		return nil, err
	}
	uidList := itemInfoList.GetSellerIDs()
	for _, item := range itemInfoList {
		userMap, err := p.userRepo.FetchByUidList(ctx, uidList)
		if err != nil {
			p.logger.Errorf("SearchItem failed to FetchByUsername, err", err)
			return out, err
		}
		userInfo := userMap[int64(item.SellerId)]
		if userInfo == nil {
			continue
		}
		itemInfo := &ItemInfoWithSeller{
			SellerID:       item.SellerId,
			SellerAvatar:   userInfo.Avatar,
			SellerNickName: userInfo.Nickname,
			ItemId:         item.ItemId,
			ItemName:       item.ItemName,
			IconUrl:        item.IconUrl,
			Price:          item.Price,
			BookedCnt:      item.BookedCnt,
		}
		out = append(out, itemInfo)
	}
	return out, nil
}
