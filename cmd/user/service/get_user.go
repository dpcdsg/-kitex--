package service

import (
	"github.com/ozline/tiktok/cmd/user/dal/db"
	"github.com/ozline/tiktok/cmd/user/pack"
	"github.com/ozline/tiktok/kitex_gen/user"
)

func (s *UserService) GetUser(req *user.InfoRequest) (*user.User, error) {
	userModel, err := db.GetUserByID(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return pack.User(userModel), nil
}
