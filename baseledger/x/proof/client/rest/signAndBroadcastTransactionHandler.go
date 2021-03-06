package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/Baseledger/baseledger/common"
	"github.com/Baseledger/baseledger/logger"
	baseledgerTypes "github.com/Baseledger/baseledger/x/proof/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type signAndBroadcastTransactionRequest struct {
	TransactionId string `json:"transaction_id"`
	Payload       string `json:"payload"`
	OpCode        uint32 `json:"op_code"`
}

func signAndBroadcastTransactionHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := parseSignAndBroadcastTransactionRequest(w, r, clientCtx)

		if !isValidUUID(req.TransactionId) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "transaction id must be uuid")
			return
		}

		clientCtx, err := BuildClientCtx(clientCtx)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		accNum, accSeq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(*clientCtx, clientCtx.FromAddress)

		if err != nil {
			logger.Errorf("error while retrieving acc %v\n", err.Error())
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "error while retrieving acc")
			return
		}

		balanceOk, err := checkTokenBalance(clientCtx.GetFromAddress().String(), req.Payload)

		if err != nil {
			logger.Errorf("check balance failed %v\n", err)
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "error while checking balance")
			return
		}

		if !balanceOk {
			logger.Errorf("check balance failed %v\n", err)
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "not enough tokens")
			return
		}

		msg := baseledgerTypes.NewMsgCreateBaseledgerTransaction(clientCtx.GetFromAddress().String(), req.TransactionId, req.Payload, req.OpCode)
		if err := msg.ValidateBasic(); err != nil {
			logger.Errorf("msg validate basic failed %v\n", err.Error())
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		logger.Infof("msg with encrypted payload to be broadcasted %s\n", msg)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		txHash, err := BroadcastAndGetTxHash(*clientCtx, msg, accNum, accSeq, false)

		if err != nil {
			logger.Errorf("broadcasting failed %v\n", err.Error())
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		logger.Infof("broadcasted tx hash %v\n", *txHash)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(*txHash))
		w.WriteHeader(http.StatusOK)
		return
	}
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func checkBalanceHandler(clientCtx client.Context, payload string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, err := BuildClientCtx(clientCtx)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		balanceOk, err := checkTokenBalance(clientCtx.GetFromAddress().String(), payload)

		if err != nil {
			logger.Errorf("check balance failed %v\n", err)
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "error while checking balance")
			return
		}

		if !balanceOk {
			logger.Errorf("check balance failed %v\n", err)
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "not enough tokens")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}
}

func checkTokenBalance(address string, payload string) (bool, error) {
	grpcConn, err := grpc.Dial(
		"127.0.0.1:9090",
		// The SDK doesn't support any transport security mechanism.
		grpc.WithInsecure(),
	)

	defer grpcConn.Close()

	if err != nil {
		logger.Errorf("grpc conn failed %v\n", err.Error())
		return false, err
	}

	queryClient := banktypes.NewQueryClient(grpcConn)
	res, err := queryClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{Address: address, Denom: common.WorkTokenDenom})

	if err != nil {
		logger.Errorf("grpc query failed %v\n", err.Error())
		return false, err
	}

	logger.Infof("found acc balance %v\n", res.Balance.Amount)

	payloadFee, err := common.CalcWorkTokenFeeBasedOnPayloadSize(payload)

	if err != nil {
		return false, errors.New("Error while calculating fee")
	}

	return res.Balance.Amount.GTE(payloadFee.AmountOf("work")), nil
}

func parseSignAndBroadcastTransactionRequest(w http.ResponseWriter, r *http.Request, clientCtx client.Context) *signAndBroadcastTransactionRequest {
	var req signAndBroadcastTransactionRequest
	if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
		return nil
	}

	return &req
}
