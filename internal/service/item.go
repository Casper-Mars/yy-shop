package service

import (
	"context"
	v1 "yy-shop/api/v1"
	"yy-shop/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type itemService struct {
	v1.UnimplementedProductServer
	log        *log.Helper
	productMgr *biz.ProductMgr
}

func NewProductServer(logger log.Logger, productMgr *biz.ProductMgr) v1.ProductServer {
	return &itemService{
		log:        log.NewHelper(logger),
		productMgr: productMgr,
	}
}

func (a *itemService) SearchItem(ctx context.Context, request *v1.SearchItemRequest) (*v1.SearchItemResponse, error) {
	out := &v1.SearchItemResponse{}
	pageToken, err := a.parseToken(ctx, request.GetPageToken())
	if err != nil {
		a.log.Errorf("SearchItem failed to parseToken, err:%v", err)
		return out, err
	}
	itemInfoList, err := a.productMgr.SearchItem(ctx, request.GetName(), pageToken, request.GetPageSize())
	if err != nil {
		a.log.Errorf("SearchItem failed, request:%v, err:%v", request, err)
		return out, err
	}

	newToken, err := a.genNewToken(ctx, request.GetPageToken(), request.GetPageSize())
	if err != nil {
		a.log.Errorf("SearchItem failed to genNewToken, err:%v", err)
	}

	return &v1.SearchItemResponse{
		ItemList:  a.convertItemInfo2Pb(itemInfoList),
		PageToken: newToken,
	}, nil
}

func (a *itemService) parseToken(ctx context.Context, token string) (uint32, error) {
	return 0, nil
}

func (a *itemService) genNewToken(ctx context.Context, token string, pageSize uint32) (string, error) {
	return "", nil
}

func (a *itemService) convertItemInfo2Pb(itemInfoList []*biz.ItemInfoWithSeller) []*v1.Item {
	return []*v1.Item{}
}
