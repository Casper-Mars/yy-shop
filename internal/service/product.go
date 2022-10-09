package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	v1 "yy-shop/api/v1"
	"yy-shop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type itemService struct {
	v1.UnimplementedProductServer
	log           *log.Helper
	productMgr    *biz.ProductMgr
	searchService biz.SearchMgr
}

var (
	ErrTokenInvalid = errors.New("Token异常")
)

type pageToken struct {
	TokenID  uint32 `json:"token_id"`  // 已获取的最大id
	ItemName string `json:"item_name"` // 商品名称
}

func (p *pageToken) string() string {
	if p == nil {
		return ""
	}
	jsonByte, err := json.Marshal(p)
	if err != nil {
		log.Errorf("pageToken failed to Marshal, err:%v", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(jsonByte)
}

func NewProductServer(logger log.Logger, productMgr *biz.ProductMgr) v1.ProductServer {
	return &itemService{
		log:        log.NewHelper(logger),
		productMgr: productMgr,
	}
}

func (a *itemService) SearchItem(ctx context.Context, request *v1.SearchItemRequest) (*v1.SearchItemResponse, error) {
	out := &v1.SearchItemResponse{}
	tokenStr := request.GetPageToken()
	token := a.parseToken(ctx, tokenStr)
	if tokenStr != "" && token.ItemName != request.Name {
		log.Errorf("SearchItem failed, err:%v", ErrTokenInvalid)
		return out, ErrTokenInvalid
	}
	result, err := a.searchService.Search(ctx, nil)
	itemIDs := result.GetAllID().AsUint32()
	itemInfoList, err := a.productMgr.SearchItem(ctx, itemIDs...)
	itemLen := len(itemInfoList)
	if err != nil || itemLen == 0 {
		a.log.Errorf("SearchItem failed, request:%v, err:%v", request, err)
		return out, err
	}

	retToken := ""
	if itemLen == int(request.GetPageSize()) {
		retToken = a.genNewToken(ctx, request.GetName(), request.GetPageToken(), itemInfoList[len(itemInfoList)-1].ItemId)
	}
	return &v1.SearchItemResponse{
		ItemList:  a.convertItemInfo2Pb(itemInfoList),
		PageToken: retToken,
	}, nil
}

func (a *itemService) parseToken(ctx context.Context, token string) *pageToken {
	out := &pageToken{}
	if token == "" {
		return out
	}

	jsonByte, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.Errorf("DecodeString failed, err:%v", err)
		return out
	}

	err = json.Unmarshal(jsonByte, out)
	if err != nil {
		log.Errorf("Unmarshal failed, err:%v", err)
		return out
	}
	return out
}

func (a *itemService) genNewToken(ctx context.Context, itemName, token string, tokenID uint32) string {
	if token == "" {
		return a.initToken(ctx, itemName, tokenID)
	}
	jsonByte, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.Errorf("DecodeString failed, err:%v", err)
		return ""
	}
	pt := &pageToken{}
	err = json.Unmarshal(jsonByte, pt)
	if err != nil {
		log.Errorf("Unmarshal failed, err:%v", err)
		return ""
	}
	pt.TokenID = tokenID
	return pt.string()
}

func (a *itemService) initToken(ctx context.Context, itemName string, token uint32) string {
	pt := &pageToken{TokenID: token, ItemName: itemName}
	return pt.string()
}
func (a *itemService) convertItemInfo2Pb(itemInfoList []*biz.ItemInfoWithSeller) []*v1.Item {
	out := make([]*v1.Item, 0, len(itemInfoList))
	for _, item := range itemInfoList {
		i := &v1.Item{
			UserInfo: &v1.Item_UserInfo{
				Uid:      item.SellerID,
				NickName: item.SellerNickName,
				Avatar:   item.SellerAvatar,
			},
			ItemInfo: &v1.Item_ItemInfo{
				ItemId:    item.ItemId,
				ItemName:  item.ItemName,
				ItemCover: item.IconUrl,
				Price:     float32(item.Price),
			},
		}
		out = append(out, i)
	}
	return out
}
