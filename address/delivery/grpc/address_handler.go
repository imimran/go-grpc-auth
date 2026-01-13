package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/imimran/go-grpc-auth/address/domain"
	transformer "github.com/imimran/go-grpc-auth/address/transformer/grpc"
	"github.com/imimran/go-grpc-auth/address/usecase"
	pb "github.com/imimran/go-grpc-auth/proto" // Ensure this matches your pb package path
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AddressHandler struct {
	addressUC *usecase.AddressUsecase
	pb.UnimplementedAddressServiceServer
}

func NewAddressHandler(addressUC *usecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{addressUC: addressUC}
}

func (h *AddressHandler) ListAddress(ctx context.Context, req *pb.AddressListRequest) (*pb.AddressListResponse, error) {
    // 1. Set default pagination values
    page := int(req.GetPage())
    if page <= 0 { page = 1 }
    
    limit := int(req.GetLimit())
    if limit <= 0 { limit = 10 }

    // 2. Call Usecase/Repository
    addresses, total, err := h.addressUC.List(ctx, page, limit)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to fetch addresses: %v", err)
    }

    // 3. Transform to Proto
    var pbAddresses []*pb.Address
    for _, addr := range addresses {
        pbAddresses = append(pbAddresses, transformer.ToProtoAddress(&addr))
    }

    // 4. Calculate total pages
    totalPages := int32((total + int64(limit) - 1) / int64(limit))

    return &pb.AddressListResponse{
        Addresses:    pbAddresses,
        TotalRecords: total,
        TotalPages:   totalPages,
    }, nil
}

func (h *AddressHandler) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.Address, error) {
	// Map request to Domain model using the nested Coordinates struct
	addr := &domain.Address{
		RawAddress: req.RawAddress,
		Coordinates: domain.Coordinates{
			Latitude:  req.Coordinates.GetLatitude(),
			Longitude: req.Coordinates.GetLongitude(),
		},
		Accuracy: req.Accuracy,
		Source:   req.Source,
	}

	err := h.addressUC.Create(ctx, addr)
	if err != nil {
		return nil, err
	}
	
	return transformer.ToProtoAddress(addr), nil
}

func (h *AddressHandler) GetAddress(ctx context.Context, req *pb.AddressId) (*pb.Address, error) {
    address, err := h.addressUC.GetByID(ctx, req.GetId())
    if err != nil {
        // Check if the error is "Record Not Found"
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, status.Error(codes.NotFound, "Address not found in database")
        }
        // Other unexpected errors remain Internal
        return nil, status.Errorf(codes.Internal, "Database error: %v", err)
    }
    
    return transformer.ToProtoAddress(address), nil
}


func (h *AddressHandler) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.Address, error) {
    // 1. Parse string ID to uuid.UUID
    parsedID, err := uuid.Parse(req.Id)
    if err != nil {
        return nil, status.Errorf(codes.InvalidArgument, "invalid uuid format: %v", err)
    }

    // 2. Map request to Domain model
    addr := &domain.Address{
        ID:         parsedID, // Crucial: GORM uses this to find the record to update
        RawAddress: req.RawAddress,
        Coordinates: domain.Coordinates{
            Latitude:  req.Coordinates.GetLatitude(),
            Longitude: req.Coordinates.GetLongitude(),
        },
        Accuracy: req.Accuracy,
        Source:   req.Source,
    }

    // 3. Call Usecase
    if err := h.addressUC.Update(ctx, addr); err != nil {
        return nil, status.Errorf(codes.Internal, "update failed: %v", err)
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