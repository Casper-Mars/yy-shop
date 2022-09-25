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
	bookedCnt      uint32  // 想要的人数
}

//
type ItemInfo struct {
	ItemId   uint32  // 商品ID
	ItemName string  // 商品名称
	Price    float64 // 价格
	IconUrl  string  // 商品图片连接
	sellerId uint32  // 卖家id
}

type ProductRepo interface {
	// FetchByUsername 获取指定用户名的用户的信息，如果用户不存在，则返回 ErrUserNotExist。
	FetchByItemName(ctx context.Context, itemName string, pageToken, pageSize uint32) (itemInfoList []*ItemInfo, err error)
}

type ProductCache interface {
	GetItemBookCnt(ctx context.Context, itemId uint32) (uint32, error)
}

type ProductMgr struct {
	userRepo     UserRepo
	producctRepo ProductRepo
	productCache ProductCache
	logger       *log.Helper
}

//NewAccountUseCase 创建一个AccountUseCase，依赖作为参数传入
func NewProductMgr(logger log.Logger, userRepo UserRepo, producctRepo ProductRepo, prodcutCache ProductCache) *ProductMgr {
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
	for _, item := range itemInfoList {
		user, err := p.userRepo.FetchByUid(ctx, int64(item.sellerId))
		if err != nil {
			p.logger.Errorf("SearchItem failed to FetchByUsername, err", err)
			return out, err
		}
		bookCnt, err := p.productCache.GetItemBookCnt(ctx, item.ItemId)
		if err != nil {
			// 获取失败默认为0
			p.logger.Errorf("SearchItem failed to GetItemBookCnt, err:%v", err)
		}
		itemInfo := &ItemInfoWithSeller{
			SellerID:       uint32(user.ID),
			SellerAvatar:   user.Avatar,
			SellerNickName: user.Nickname,
			ItemId:         item.ItemId,
			ItemName:       item.ItemName,
			IconUrl:        item.IconUrl,
			Price:          item.Price,
			bookedCnt:      bookCnt,
		}
		out = append(out, itemInfo)
	}
	return out, nil
}
