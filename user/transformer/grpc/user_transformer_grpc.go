package transformer

import (
	pb "github.com/imimran/go-grpc-auth/proto"
	"github.com/imimran/go-grpc-auth/user/domain"
)

// ToProtoUser converts a domain.User model to a gRPC pb.User message.
// It explicitly omits the hashed password for security.
func ToProtoUser(user *domain.User) *pb.User {
	return &pb.User{
		Id:    user.ID,
		Email: user.Email,
		FullName: user.FullName,
	}
}

// ToProtoUserList converts a slice of domain.User to a slice of pb.User.
func ToProtoUserList(users []*domain.User) []*pb.User {
	protoUsers := make([]*pb.User, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, ToProtoUser(user))
	}
	return protoUsers
}
