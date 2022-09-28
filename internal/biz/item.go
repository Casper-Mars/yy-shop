package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrItemNotExist = errors.New("商品名称有误")
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

//
type ItemInfo struct {
	ItemId    uint32  `gorm:"column:id"`         // 商品ID
	ItemName  string  `gorm:"column:item_name"`  // 商品名称
	IconUrl   string  `gorm:"column:icon_url"`   // 商品图片连接
	Price     float64 `gorm:"column:price"`      // 价格
	SellerId  uint32  `gorm:"column:seller_id"`  // 卖家id
	BookedCnt uint32  `gorm:"column:booked_cnt"` // 想要的人数
}

type ItemRepo interface {
	// FetchByUsername 获取指定用户名的用户的信息，如果用户不存在，则返回 ErrUserNotExist。
	FetchByItemName(ctx context.Context, itemName string, pageToken, pageSize uint32) (itemInfoList []*ItemInfo, err error)
}

type ProductCache interface {
	GetItemBookCnt(ctx context.Context, itemId uint32) (uint32, error)
}

type ProductMgr struct {
	userRepo     UserRepo
	producctRepo ItemRepo
	productCache ProductCache
	logger       *log.Helper
}

//NewAccountUseCase 创建一个AccountUseCase，依赖作为参数传入
func NewProductMgr(logger log.Logger, userRepo UserRepo, producctRepo ItemRepo, prodcutCache ProductCache) *ProductMgr {
	return &ProductMgr{
		userRepo:     userRepo,
		producctRepo: producctRepo,
		productCache: prodcutCache,
		logger:       log.NewHelper(logger),
	}
}

//Register 注册
func (p *ProductMgr) SearchItem(ctx context.Context, itemName string, pageToken, pageSize uint32) ([]*ItemInfoWithSeller, error) {
	out := make([]*ItemInfoWithSeller, 0)
	itemInfoList, err := p.producctRepo.FetchByItemName(ctx, itemName, pageToken, pageSize)
	if err != nil {
		p.logger.Errorf("SearchItem failed to FetchByItemName, err:%v", err)
		return out, err
	}
	if len(itemInfoList) == 0 {
		log.Errorf("SearchItem item do not exist, itemName:%s", itemName)
		return out, ErrItemNotExist
	}
	uidList := make([]int64, 0, len(itemInfoList))
	for _, item := range itemInfoList {
		uidList = append(uidList, int64(item.SellerId))
	}
	for _, item := range itemInfoList {
		userMap, err := p.userRepo.FetchByUidList(ctx, uidList)
		if err != nil {
			p.logger.Errorf("SearchItem failed to FetchByUsername, err", err)
			return out, err
		}
		itemInfo := &ItemInfoWithSeller{
			SellerID:       uint32(item.SellerId),
			SellerAvatar:   userMap[int64(item.SellerId)].Avatar,
			SellerNickName: userMap[int64(item.SellerId)].Nickname,
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
