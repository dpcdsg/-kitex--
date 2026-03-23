package pack

import (
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/kitex_gen/user"
)

func strPtr(s string) *string {
	return &s
}

func User(u *user.User) *api.User {
	if u == nil {
		return nil
	}
	return &api.User{
		ID:        u.Id,
		Name:      u.Name,
		Avatar:    strPtr(u.Avatar),
		Signature: strPtr(u.Signature),
	}
}
