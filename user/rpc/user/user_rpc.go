package user

import (
	"context"
	"projects/LDmitryLD/task-service/user/internal/models"
	"projects/LDmitryLD/task-service/user/internal/modules/user/service"

	pb "github.com/LDmitryLD/protos2/gen/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceGRPC struct {
	userService service.Userer
	pb.UnimplementedUsererServer
}

func NewUserServiceGRPC(userService service.Userer) *UserServiceGRPC {
	return &UserServiceGRPC{userService: userService}
}

func (u *UserServiceGRPC) Profile(ctx context.Context, in *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	user, err := u.userService.Profile(ctx, int(in.GetId()))
	if err != nil {
		return nil, err
	}

	return &pb.ProfileResponse{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Tasks: user.Tasks}, nil
}

func (u *UserServiceGRPC) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := u.userService.Create(models.User{FirstName: in.FirstName, LastName: in.LastName, Email: in.Email})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{Id: uint32(id)}, nil
}

func (u *UserServiceGRPC) Exists(ctx context.Context, in *pb.ExistsRequest) (*pb.ExistsResponse, error) {
	err := u.userService.Exists(ctx, int(in.Id))
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return &pb.ExistsResponse{Exists: false}, nil
		default:
			return nil, err
		}
	}

	return &pb.ExistsResponse{Exists: true}, nil
}
