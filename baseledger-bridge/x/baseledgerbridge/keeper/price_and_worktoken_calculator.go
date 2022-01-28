package keeper

import (
	"math/big"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalculateAvgUbtPriceForAttestation(att types.Attestation) *big.Int {

	sum := big.NewInt(0)

	for i := 0; i < len(att.UbtPrices); i++ {
		sum.Add(sum, att.UbtPrices[i].BigInt())
	}

	arrayLengthBigInt := big.NewInt(int64(len(att.UbtPrices)))

	mean := sum.Div(sum, arrayLengthBigInt)

	standardDev := big.NewInt(0)

	for j := 0; j < len(att.UbtPrices); j++ {

		diffFromMean := new(big.Int).Sub(att.UbtPrices[j].BigInt(), mean)

		diffFromMeanSquared := diffFromMean.Exp(diffFromMean, big.NewInt(2), nil)

		standardDev.Add(standardDev, diffFromMeanSquared)
	}

	standardDev.Sqrt(new(big.Int).Div(standardDev, arrayLengthBigInt))

	var cleansedPriceArray []big.Int
	for k := 0; k < len(att.UbtPrices); k++ {
		ubtPrice := att.UbtPrices[k].BigInt()

		oneStDevLessFromMean := new(big.Int).Sub(mean, standardDev)
		oneStDevGreaterFromMean := new(big.Int).Add(mean, standardDev)

		if (ubtPrice.Cmp(oneStDevLessFromMean) == +1 || ubtPrice.Cmp(oneStDevLessFromMean) == 0) &&
			(ubtPrice.Cmp(oneStDevGreaterFromMean) == -1 || ubtPrice.Cmp(oneStDevGreaterFromMean) == 0) {
			cleansedPriceArray = append(cleansedPriceArray, *ubtPrice)
		}
	}

	cleansedSum := big.NewInt(0)

	for l := 0; l < len(cleansedPriceArray); l++ {
		cleansedSum.Add(cleansedSum, &cleansedPriceArray[l])
	}

	return cleansedSum.Div(cleansedSum, big.NewInt(int64(len(cleansedPriceArray))))
}

func CalculateAmountOfWorkTokens(depositedUbtAmount *big.Int, averagePrice *big.Int) *big.Int {
	// TODO: BAS-121 - Move this hardcoded value to config or somewhere
	// TODO: Ognjen - Verify calculation
	worktokenEurPrice, _ := sdk.NewDecFromStr("0.1")
	worktokenEurPriceInt := worktokenEurPrice.BigInt()

	depositedEurValueInt := depositedUbtAmount.Mul(depositedUbtAmount, averagePrice)

	return depositedEurValueInt.Div(depositedEurValueInt, worktokenEurPriceInt)
}
