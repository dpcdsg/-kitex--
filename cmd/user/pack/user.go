package pack

import (
	"github.com/ozline/tiktok/kitex_gen/user"

	"github.com/ozline/tiktok/cmd/user/dal/db"
)

func User(data *db.User) *user.User {
	if data == nil {
		return nil
	}

	return &user.User{
		Id:        data.Id,
		Name:      data.Username,
		Avatar:    data.Avatar,
		Signature: data.Signature,
		CreatedAt: data.CreatedAt.Format("2006-01-02 15:05"),
	}
}
