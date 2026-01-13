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
		Id:                addr.ID.String(), // Convert uuid.UUID to string
		RawAddress:        addr.RawAddress,
		NormalizedAddress: addr.NormalizedAddress,
		Latitude:          addr.Latitude,
		Longitude:         addr.Longitude,
		Accuracy:          addr.Accuracy,
		Source:            addr.Source,
	}
}

// ToProtoAddressList converts a slice of Domain Addresses to a slice of Protobuf Addresses
func ToProtoAddressList(addresses []*domain.Address) []*pb.Address {
	protoAddresses := make([]*pb.Address, len(addresses))
	
	for i, addr := range addresses {
		protoAddresses[i] = ToProtoAddress(addr)
	}
	
	return protoAddresses
}