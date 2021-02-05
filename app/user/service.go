package user

import (
	"context"

	pb "github.com/DoNewsCode/skeleton/app/proto"
	"github.com/DoNewsCode/skeleton/internal/entities"
	"github.com/pkg/errors"
)

type UserRepository interface {
	Find(ctx context.Context, id uint) (*entities.User, error)
}

type Service struct {
	UserRepository UserRepository
}

func (s Service) Login(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
	user, err := s.UserRepository.Find(ctx, uint(in.Id))
	if err != nil {
		return nil, errors.Wrap(err, "login failed")
	}

	var resp pb.UserInfoReply
	resp = pb.UserInfoReply{
		Id:        int64(user.ID),
		Name:      user.UserName,
		CreatedAt: user.CreatedAt.String(),
	}
	return &resp, nil
}
