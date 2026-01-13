package grpc

import (
	"github.com/imimran/go-grpc-auth/address/domain"
	pb "github.com/imimran/go-grpc-auth/proto"
)

// ToProtoAddress converts a single Domain Address to a Protobuf Address
func ToProtoAddress(addr *domain.Address) *pb.Address {
	if addr == nil {
		return nil
	}

	return &pb.Address{
        Id:                addr.ID.String(),
        RawAddress:        addr.RawAddress,
        NormalizedAddress: addr.NormalizedAddress,
        // Mapping nested struct to nested Proto message
        Coordinates: &pb.Coordinates{
            Latitude:  addr.Coordinates.Latitude,
            Longitude: addr.Coordinates.Longitude,
        },
        Accuracy: addr.Accuracy,
        Source:   addr.Source,
    }
}

// ToProtoAddressList converts a slice of Domain Addresses to a slice of Protobuf Addresses
// func ToProtoAddress(addr *domain.Address) *pb.Address {
//     if addr == nil {
//         return nil
//     }
//     return &pb.Address{
//         Id:                addr.ID.String(),
//         RawAddress:        addr.RawAddress,
//         NormalizedAddress: addr.NormalizedAddress,
//         Coordinates: &pb.Coordinates{
//             Latitude:  addr.Coordinates.Latitude,
//             Longitude: addr.Coordinates.Longitude,
//         },
//         Accuracy: addr.Accuracy,
//         Source:   addr.Source,
//     }
// }