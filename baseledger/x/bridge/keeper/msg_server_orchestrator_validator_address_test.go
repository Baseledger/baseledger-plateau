package keeper_test

import (
    "strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

    keepertest "github.com/Baseledger/baseledger/testutil/keeper"
    "github.com/Baseledger/baseledger/x/bridge/keeper"
    "github.com/Baseledger/baseledger/x/bridge/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestOrchestratorValidatorAddressMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.BridgeKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	validatorAddress := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateOrchestratorValidatorAddress{ValidatorAddress: validatorAddress,
		    OrchestratorAddress: strconv.Itoa(i),
            
		}
		_, err := srv.CreateOrchestratorValidatorAddress(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetOrchestratorValidatorAddress(ctx,
		    expected.OrchestratorAddress,
            
		)
		require.True(t, found)
		require.Equal(t, expected.ValidatorAddress, rst.ValidatorAddress)
	}
}

func TestOrchestratorValidatorAddressMsgServerUpdate(t *testing.T) {
	validatorAddress := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateOrchestratorValidatorAddress
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateOrchestratorValidatorAddress{ValidatorAddress: validatorAddress,
			    OrchestratorAddress: strconv.Itoa(0),
                
			},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateOrchestratorValidatorAddress{ValidatorAddress: "B",
			    OrchestratorAddress: strconv.Itoa(0),
                
			},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgUpdateOrchestratorValidatorAddress{ValidatorAddress: validatorAddress,
			    OrchestratorAddress: strconv.Itoa(100000),
                
			},
			err:     sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.BridgeKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateOrchestratorValidatorAddress{ValidatorAddress: validatorAddress,
			    OrchestratorAddress: strconv.Itoa(0),
                
			}
			_, err := srv.CreateOrchestratorValidatorAddress(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateOrchestratorValidatorAddress(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetOrchestratorValidatorAddress(ctx,
				    expected.OrchestratorAddress,
                    
				)
				require.True(t, found)
				require.Equal(t, expected.ValidatorAddress, rst.ValidatorAddress)
			}
		})
	}
}

func TestOrchestratorValidatorAddressMsgServerDelete(t *testing.T) {
	validatorAddress := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteOrchestratorValidatorAddress
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeleteOrchestratorValidatorAddress{ValidatorAddress: validatorAddress,
			    OrchestratorAddress: strconv.Itoa(0),
                
			},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeleteOrchestratorValidatorAddress{ValidatorAddress: "B",
			    OrchestratorAddress: strconv.Itoa(0),
                
			},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgDeleteOrchestratorValidatorAddress{ValidatorAddress: validatorAddress,
			    OrchestratorAddress: strconv.Itoa(100000),
                
			},
			err:     sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.BridgeKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateOrchestratorValidatorAddress(wctx, &types.MsgCreateOrchestratorValidatorAddress{ValidatorAddress: validatorAddress,
			    OrchestratorAddress: strconv.Itoa(0),
                
			})
			require.NoError(t, err)
			_, err = srv.DeleteOrchestratorValidatorAddress(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetOrchestratorValidatorAddress(ctx,
				    tc.request.OrchestratorAddress,
                    
				)
				require.False(t, found)
			}
		})
	}
}
