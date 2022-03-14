package rest

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Baseledger/baseledger/logger"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
)

const (
	errCodeMismatch = 32
)

var (
	// errors are of the form:
	// "account sequence mismatch, expected 25, got 27: incorrect account sequence"
	recoverRegexp = regexp.MustCompile(`^account sequence mismatch, expected (\d+), got (\d+):`)
)

func BuildClientCtx(clientCtx client.Context) (*client.Context, error) {
	keyringInstance, err := NewKeyringInstance()
	if err != nil {
		logger.Errorf("error getting keyring instance %v\n", err.Error())
		return nil, err
	}

	// node can specify key address as env variable, otherwise it will use first from keylist
	key, err := getKey(keyringInstance)
	if err != nil {
		logger.Errorf("error getting key %v\n", err.Error())
		return nil, err
	}

	clientCtx = clientCtx.
		WithKeyring(keyringInstance).
		WithFromAddress(key.GetAddress()).
		WithSkipConfirmation(true).
		WithFromName(key.GetName()).
		WithBroadcastMode("sync")

	return &clientCtx, nil
}

func getKey(keyringInstance keyring.Keyring) (keyring.Info, error) {
	keysList, err := keyringInstance.List()
	if err != nil {
		logger.Errorf("error getting key list %v\n", err.Error())
		return nil, errors.New("")
	}
	var key keyring.Info
	addrString := viper.GetString("KEY_ADDRESS")
	if addrString != "" {
		logger.Infof("key specified in env variables %v\n", addrString)
		addr, err := sdk.AccAddressFromBech32(addrString)
		if err != nil {
			logger.Errorf("key specified in wrong format %v\n", err.Error())
			return nil, err
		}

		key, err = keyringInstance.KeyByAddress(addr)
		if err != nil {
			logger.Errorf("error getting specified key %v\n", err.Error())
			return nil, err
		}
	} else {
		logger.Infof("key not specified in env variables, picking first from list")
		key = keysList[0]
	}

	if key == nil {
		logger.Error("key is nil")
		return nil, errors.New("key is nil")
	}
	return key, nil
}

func NewKeyringInstance() (keyring.Keyring, error) {
	input := &bytes.Buffer{}
	fmt.Printf("KEYRING KEYRING KEYRING %v %v\n", viper.GetString("KEYRING_PASSWORD"), viper.GetString("KEYRING_DIR"))
	fmt.Fprintln(input, viper.GetString("KEYRING_PASSWORD"), viper.GetString("KEYRING_PASSWORD"))
	kr, err := keyring.New("baseledger", "file", viper.GetString("KEYRING_DIR"), input)

	if err != nil {
		logger.Errorf("error fetching keyring, check if you configured KEYRING_PASSWORD and KEYRING_DIR %v\n", err.Error())
		return nil, errors.New("error fetching key ring")
	}

	return kr, nil
}

func SignTxAndGetTxBytes(clientCtx client.Context, msg sdk.Msg, accNum uint64, accSeq uint64) ([]byte, error) {
	logger.Infof("retrieved account %v %v\n", accNum, accSeq)
	txFactory := tx.Factory{}.
		WithChainID("baseledger").
		WithGas(1000000). // hardcoding gasWanted to high number since fees will allways be 1 token
		WithTxConfig(clientCtx.TxConfig).
		WithAccountNumber(accNum).
		WithSequence(accSeq).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithKeybase(clientCtx.Keyring)

	txFactory, err := prepareFactory(clientCtx, txFactory)
	if err != nil {
		logger.Errorf("prepare factory error %v\n", err.Error())
		return nil, errors.New("sign tx error")
	}
	simResp, _, err := tx.CalculateGas(clientCtx, txFactory, msg)
	if err != nil {
		logger.Errorf("calc gas error %v\n", err.Error())
		return nil, err
	}

	txFactory = txFactory.WithGas(simResp.GasInfo.GasUsed)

	transaction, err := tx.BuildUnsignedTx(txFactory, msg)
	if err != nil {
		logger.Errorf("build unsigned tx error %v\n", err.Error())
		return nil, errors.New("sign tx error")
	}

	err = tx.Sign(txFactory, clientCtx.GetFromName(), transaction, false)
	if err != nil {
		logger.Errorf("sign tx error %v\n", err.Error())
		return nil, errors.New("sign tx error")
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(transaction.GetTx())
	if err != nil {
		logger.Errorf("tx encoder %v\n", err.Error())
		return nil, errors.New("sign tx error")
	}

	return txBytes, nil
}

// copied from cosmos-sdk/client because it is not public anymore
// prepareFactory ensures the account defined by ctx.GetFromAddress() exists and
// if the account number and/or the account sequence number are zero (not set),
// they will be queried for and set on the provided Factory. A new Factory with
// the updated fields will be returned.
func prepareFactory(clientCtx client.Context, txf tx.Factory) (tx.Factory, error) {
	from := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	return txf, nil
}

func BroadcastAndGetTxHash(clientCtx client.Context, msg sdk.Msg, accNum uint64, accSeq uint64, retried bool) (*string, error) {
	txBytes, err := SignTxAndGetTxBytes(clientCtx, msg, accNum, accSeq)
	if err != nil {
		return nil, err
	}
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		logger.Errorf("error while broadcasting tx %v\n", err.Error())
		return nil, err
	}

	// if broadcast was successful, return txHash
	if res.Code == 0 {
		logMsg := "BROADCASTED"
		if retried {
			logMsg = "REBROADCASTED"
		}
		logger.Infof("TRANSACTION %v WITH RESULT %v\n", logMsg, res)
		return &res.TxHash, nil
	}

	if res.Code != 0 && res.Code != errCodeMismatch {
		logger.Errorf("broadcast failed with code different than missmatch %v %v\n", res.Code, res)
		return nil, err
	}

	// if code is missmatch and it is retrying, don't handle it and return
	if retried {
		return nil, errors.New("broadcast failed after retrying")
	}

	// if code is missmatch first time, parse log and try again
	logger.Infof("ACCOUNT SEQUENCE MISSMATCH %v\n", res.RawLog)

	nextSequence, ok := parseNextSequence(accSeq, res.RawLog)

	if !ok {
		return nil, errors.New("broadcast failed when parsing sequence")
	}

	logger.Infof("RETRYING WITH SEQUENCE %v\n", nextSequence)

	return BroadcastAndGetTxHash(clientCtx, msg, accNum, nextSequence, true)
}

func parseNextSequence(current uint64, message string) (uint64, bool) {
	// "account sequence mismatch, expected 25, got 27: incorrect account sequence"
	matches := recoverRegexp.FindStringSubmatch(message)

	if len(matches) != 3 {
		return 0, false
	}

	if len(matches[1]) == 0 || len(matches[2]) == 0 {
		return 0, false
	}

	expected, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil || expected == 0 {
		return 0, false
	}

	received, err := strconv.ParseUint(matches[2], 10, 64)
	if err != nil || received == 0 {
		return 0, false
	}

	if received != current {
		return expected, true
	}

	return expected, true
}
