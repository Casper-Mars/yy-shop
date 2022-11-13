package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
)

var (
	ErrItemNotExist = errors.New("商品名称有误")
	ErrNoResult     = errors.New("没有搜索结果")
)

type Seller struct {
	ID       uint32 // 卖家ID
	Avatar   string // 卖家头像
	NickName string // 卖家昵称
}

type ProductIndex struct {
	ID       uint32  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	SellerID uint32  `json:"seller_id"`
}

type Product struct {
	ItemId    uint32   // 商品ID
	ItemName  string   // 商品名称
	Price     float64  // 价格
	IconUrl   string   // 商品图片连接
	BookedCnt uint32   // 想要的人数
	Images    []string // 商品图片
	Seller    *Seller  // 卖家
}

type ProductList []*Product

// ItemInfo 商品信息
type ItemInfo struct {
	ItemId    uint32  `gorm:"column:id"`         // 商品ID
	ItemName  string  `gorm:"column:item_name"`  // 商品名称
	IconUrl   string  `gorm:"column:icon_url"`   // 商品图片连接
	Price     float64 `gorm:"column:price"`      // 价格
	SellerId  uint32  `gorm:"column:seller_id"`  // 卖家id
	BookedCnt uint32  `gorm:"column:booked_cnt"` // 想要的人数
	Images    string  `gorm:"column:images"`     // 商品图片, 逗号分隔
}

type ItemList []*ItemInfo

func (i ItemList) GetSellerIDs() []uint32 {
	if len(i) == 0 {
		return []uint32{}
	}
	ids := make([]uint32, len(i))
	for i, item := range i {
		ids[i] = item.SellerId
	}
	return ids
}

type ItemRepo interface {
	// FetchByItemName 获取指定用户名的用户的信息，如果用户不存在，则返回 ErrUserNotExist。
	FetchByItemName(ctx context.Context, itemName string, pageToken, pageSize uint32) (itemInfoList []*ItemInfo, err error)
	// FetchByIds 批量获取指定id的商品信息
	FetchByIds(ctx context.Context, ids ...uint32) (itemInfoList ItemList, err error)
	// Save 保存商品信息
	Save(ctx context.Context, item *ItemInfo) (uint32, error)
}

type ProductCache interface {
}

type ProductUploadReq struct {
	userId   uint32
	itemName string   // 商品名称
	price    float64  // 价格
	iconUrl  string   // 商品图片连接
	images   []string // 商品图片
}

func NewProductUploadReq(
	userId uint32,
	itemName string,
	price string,
	iconUrl string,
	images []string,
) (*ProductUploadReq, error) {
	if userId == 0 {
		return nil, fmt.Errorf("用户id不能为空")
	}
	if itemName == "" {
		return nil, fmt.Errorf("商品名称不能为空")
	}
	p := &ProductUploadReq{
		userId:   userId,
		itemName: itemName,
		iconUrl:  iconUrl,
		images:   images,
	}
	float, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return nil, fmt.Errorf("价格格式错误: %s", price)
	}
	p.price = float
	return p, nil
}

type ProductSearchReq struct {
}

type ProductUseCase struct {
	userRepo     UserRepo
	itemRepo     ItemRepo
	logger       *log.Helper
	esSearchRepo EsSearchRepo
}

//NewProductMgr 创建一个AccountUseCase，依赖作为参数传入
func NewProductMgr(logger log.Logger, userRepo UserRepo, producctRepo ItemRepo, esRepo EsSearchRepo) *ProductUseCase {
	return &ProductUseCase{
		userRepo:     userRepo,
		itemRepo:     producctRepo,
		esSearchRepo: esRepo,
		logger:       log.NewHelper(logger),
	}
}

func (p *ProductUseCase) Upload(ctx context.Context, req *ProductUploadReq) (uint32, error) {
	id, err := p.itemRepo.Save(ctx, &ItemInfo{
		ItemName: req.itemName,
		IconUrl:  req.iconUrl,
		Price:    req.price,
		SellerId: req.userId,
		Images:   strings.Join(req.images, ","),
	})
	if err != nil {
		return 0, fmt.Errorf("保存商品信息失败: %w", err)
	}
	err = p.esSearchRepo.Upsert(ctx, "product", strconv.FormatUint(uint64(id), 10), &ProductIndex{
		ID:       id,
		Name:     req.itemName,
		Price:    req.price,
		SellerID: req.userId,
	})
	if err != nil {
		p.logger.Errorf("保存商品信息到es失败: %s", err)
	}
	return id, nil
}

//SearchItem 搜索商品
func (p *ProductUseCase) SearchItem(ctx context.Context, ids ...uint32) ([]*Product, error) {

	out := make([]*Product, 0)
	itemInfoList, err := p.itemRepo.FetchByIds(ctx, ids...)
	if err != nil {
		p.logger.Errorf("SearchItem failed to FetchByItemName, err:%v", err)
		return nil, err
	}
	// 获取卖家信息
	uidList := itemInfoList.GetSellerIDs()
	for _, item := range itemInfoList {
		userMap, err := p.userRepo.FetchByUidList(ctx, uidList)
		if err != nil {
			p.logger.Errorf("SearchItem failed to FetchByUsername, err", err)
			return out, err
		}
		userInfo := userMap[item.SellerId]
		if userInfo == nil {
			continue
		}
		itemInfo := &Product{
			ItemId:    item.ItemId,
			ItemName:  item.ItemName,
			IconUrl:   item.IconUrl,
			Price:     item.Price,
			BookedCnt: item.BookedCnt,
			Seller: &Seller{
				ID:       item.SellerId,
				Avatar:   userInfo.Avatar,
				NickName: userInfo.Nickname,
			},
		}
		out = append(out, itemInfo)
	}
	return out, nil
}
