package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	v1 "yy-shop/api/v1"
	"yy-shop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type itemService struct {
	v1.UnimplementedProductServer
	log        *log.Helper
	productMgr *biz.ProductMgr
}

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

	itemInfoList, err := a.productMgr.SearchItem(ctx, request.GetName(), a.parseToken(ctx, request.GetPageToken()),
		request.GetPageSize())
	if err != nil || len(itemInfoList) == 0 {
		a.log.Errorf("SearchItem failed, request:%v, err:%v", request, err)
		return out, err
	}
	return &v1.SearchItemResponse{
		ItemList:  a.convertItemInfo2Pb(itemInfoList),
		PageToken: a.genNewToken(ctx, request.GetName(), request.GetPageToken(), itemInfoList[len(itemInfoList)-1].ItemId),
	}, nil
}

func (a *itemService) parseToken(ctx context.Context, token string) uint32 {
	if token == "" {
		return 0
	}

	jsonByte, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.Errorf("DecodeString failed, err:%v", err)
		return 0
	}
	pt := &pageToken{}
	err = json.Unmarshal(jsonByte, pt)
	if err != nil {
		log.Errorf("Unmarshal failed, err:%v", err)
		return 0
	}
	return pt.TokenID
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
	return []*v1.Item{}
}
