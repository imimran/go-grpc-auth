package grpc

import (
	"context"

	"github.com/imimran/go-grpc-auth/address/domain"
	"github.com/imimran/go-grpc-auth/address/usecase"
	transformer "github.com/imimran/go-grpc-auth/address/transformer/grpc"
	pb "github.com/imimran/go-grpc-auth/proto"
)

type AddressHandler struct {
	addressUC *usecase.AddressUsecase
	pb.UnimplementedAddressServiceServer
}

func NewAddressHandler(addressUC *usecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{addressUC: addressUC}
}

func (h *AddressHandler) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.Address, error) {
	// Map request to Domain model
	addr := &domain.Address{
		RawAddress: req.RawAddress,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Accuracy:   req.Accuracy,
		Source:     req.Source,
	}

	err := h.addressUC.Create(ctx, addr)
	if err != nil {
		return nil, err
	}
	
	return transformer.ToProtoAddress(addr), nil
}

func (h *AddressHandler) GetAddress(ctx context.Context, req *pb.AddressId) (*pb.Address, error) {
	address, err := h.addressUC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	
	return transformer.ToProtoAddress(address), nil
}

func (h *AddressHandler) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.Address, error) {
	// Reconstruct domain object for update
	addr := &domain.Address{
		RawAddress: req.RawAddress,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Accuracy:   req.Accuracy,
		Source:     req.Source,
	}

	err := h.addressUC.Update(ctx, addr)
	if err != nil {
		return nil, err
	}
	
	return transformer.ToProtoAddress(addr), nil
}

func (h *AddressHandler) DeleteAddress(ctx context.Context, req *pb.AddressId) (*pb.Empty, error) {
	err := h.addressUC.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	
	return &pb.Empty{}, nil
}