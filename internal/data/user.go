package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"yy-shop/internal/biz"
)

type userRepo struct {
	data  *Data
	log   *log.Helper
	table *gorm.DB
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data:  data,
		log:   log.NewHelper(logger),
		table: data.db.Table("user"),
	}
}

func (u *userRepo) FetchByUsername(ctx context.Context, username string) (user *biz.User, err error) {
	user = &biz.User{}
	u.table.WithContext(ctx).First(user, "username = ?", username)
	if user.ID == 0 {
		return nil, biz.ErrUserNotExist
	}
	return user, nil
}

func (u *userRepo) Save(ctx context.Context, user *biz.User) (id int64, err error) {
	result := u.table.WithContext(ctx).Create(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}
