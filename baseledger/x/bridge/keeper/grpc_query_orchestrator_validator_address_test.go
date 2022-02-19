package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/Baseledger/baseledger/testutil/keeper"
	"github.com/Baseledger/baseledger/testutil/nullify"
	"github.com/Baseledger/baseledger/x/bridge/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestOrchestratorValidatorAddressQuerySingle(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t, true)
	keeper := testKeepers.BridgeKeeper
	ctx := testKeepers.Context
	wctx := sdk.WrapSDKContext(ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetOrchestratorValidatorAddressRequest
		response *types.QueryGetOrchestratorValidatorAddressResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetOrchestratorValidatorAddressRequest{
				OrchestratorAddress: keepertest.OrchAddrs[0].String(),
			},
			response: &types.QueryGetOrchestratorValidatorAddressResponse{OrchestratorValidatorAddress: types.OrchestratorValidatorAddress{
				ValidatorAddress:    keepertest.ValAddrs[0].String(),
				OrchestratorAddress: keepertest.OrchAddrs[0].String(),
			}},
		},
		{
			desc: "Second",
			request: &types.QueryGetOrchestratorValidatorAddressRequest{
				OrchestratorAddress: keepertest.OrchAddrs[1].String(),
			},
			response: &types.QueryGetOrchestratorValidatorAddressResponse{OrchestratorValidatorAddress: types.OrchestratorValidatorAddress{
				ValidatorAddress:    keepertest.ValAddrs[1].String(),
				OrchestratorAddress: keepertest.OrchAddrs[1].String(),
			}},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetOrchestratorValidatorAddressRequest{
				OrchestratorAddress: strconv.Itoa(100000),
			},
			err: status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.OrchestratorValidatorAddress(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestOrchestratorValidatorAddressQueryPaginated(t *testing.T) {
	testKeepers := keepertest.SetFiveValidators(t, true)
	keeper := testKeepers.BridgeKeeper
	ctx := testKeepers.Context
	wctx := sdk.WrapSDKContext(ctx)

	msgs := [5]types.OrchestratorValidatorAddress{}
	for i := range []int{0, 1, 2, 3, 4} {
		msgs[i] = types.OrchestratorValidatorAddress{
			ValidatorAddress:    keepertest.ValAddrs[i].String(),
			OrchestratorAddress: keepertest.OrchAddrs[i].String(),
		}
	}

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllOrchestratorValidatorAddressRequest {
		return &types.QueryAllOrchestratorValidatorAddressRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.OrchestratorValidatorAddressAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.OrchestratorValidatorAddress), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.OrchestratorValidatorAddress),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.OrchestratorValidatorAddressAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.OrchestratorValidatorAddress), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.OrchestratorValidatorAddress),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.OrchestratorValidatorAddressAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.OrchestratorValidatorAddress),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.OrchestratorValidatorAddressAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
