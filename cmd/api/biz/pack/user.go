package pack

import (
	api "github.com/ozline/tiktok/cmd/api/biz/model/api"
	"github.com/ozline/tiktok/kitex_gen/user"
)

func User(u *user.User) *api.User {
	if u == nil {
		return nil
	}
	out := &api.User{
		ID:   u.Id,
		Name: u.Name,
	}
	if u.Avatar != "" {
		a := u.Avatar
		out.Avatar = &a
	}
	if u.Signature != "" {
		s := u.Signature
		out.Signature = &s
	}
	return out
}
