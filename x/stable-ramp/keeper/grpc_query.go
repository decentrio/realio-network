package keeper

import (
	"context"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkquery "github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/realiotech/realio-network/x/stable-ramp/types"
)

var _ types.QueryServer = queryServer{}

func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (q queryServer) Package(
	ctx context.Context,
	req *types.QueryPackageRequest,
) (*types.QueryPackageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	key := collections.Join(req.Nonce, req.ConnectionId)

	// Try deposit packages first, then withdraw.
	if pkg, err := q.k.DepositPackages.Get(ctx, key); err == nil {
		return &types.QueryPackageResponse{Package: pkg}, nil
	}
	if pkg, err := q.k.WithdrawPackages.Get(ctx, key); err == nil {
		return &types.QueryPackageResponse{Package: pkg}, nil
	}

	return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "package not found")
}

func (q queryServer) PendingPackages(
	ctx context.Context,
	req *types.QueryPendingPackagesRequest,
) (*types.QueryPendingPackagesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	pageReq := req.Pagination
	if pageReq == nil {
		pageReq = &sdkquery.PageRequest{}
	}
	limit := pageReq.Limit
	if limit == 0 {
		limit = sdkquery.DefaultLimit
	}

	var packages []types.TrackedPackage
	var total uint64

	collect := func(_ collections.Pair[uint64, string], pkg types.TrackedPackage) (bool, error) {
		if pkg.PackageData == nil || pkg.PackageData.Status != types.PackageStatus_PACKAGE_STATUS_PENDING {
			return false, nil
		}
		if total >= pageReq.Offset && uint64(len(packages)) < limit {
			packages = append(packages, pkg)
		}
		total++
		return false, nil
	}

	if err := q.k.DepositPackages.Walk(ctx, nil, collect); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := q.k.WithdrawPackages.Walk(ctx, nil, collect); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pageRes := &sdkquery.PageResponse{}
	if pageReq.CountTotal {
		pageRes.Total = total
	}

	return &types.QueryPendingPackagesResponse{
		Packages:   packages,
		Pagination: pageRes,
	}, nil
}
