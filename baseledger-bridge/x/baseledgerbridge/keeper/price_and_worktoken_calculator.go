package keeper

import (
	"math/big"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalculateAvgUbtPriceForAttestation(att types.Attestation) *big.Int {

	mean := calcMean(att.UbtPrices)
	sd := calcStandardDeviation(att.UbtPrices, mean)

	cleansedPriceArray := getPricesArrayWithoutOutliers(att.UbtPrices, mean, sd)

	return calcAvgPrice(cleansedPriceArray)
}

func CalculateAmountOfWorkTokens(depositedUbtAmount *big.Int, averagePrice *big.Int) *big.Int {
	// TODO: BAS-121 - Move this hardcoded value to config or somewhere
	// TODO: Ognjen - Verify calculation
	worktokenEurPrice, _ := sdk.NewDecFromStr("0.1")
	worktokenEurPriceInt := worktokenEurPrice.BigInt()

	depositedEurValueInt := depositedUbtAmount.Mul(depositedUbtAmount, averagePrice)

	return depositedEurValueInt.Div(depositedEurValueInt, worktokenEurPriceInt)
}

func calcMean(prices []sdk.Int) *big.Int {
	sum := big.NewInt(0)

	for i := 0; i < len(prices); i++ {
		sum.Add(sum, prices[i].BigInt())
	}

	arrayLengthBigInt := big.NewInt(int64(len(prices)))

	return new(big.Int).Div(sum, arrayLengthBigInt)
}

func calcStandardDeviation(prices []sdk.Int, mean *big.Int) *big.Int {
	arrayLengthBigInt := big.NewInt(int64(len(prices)))
	standardDev := big.NewInt(0)

	for i := 0; i < len(prices); i++ {

		diffFromMean := new(big.Int).Sub(prices[i].BigInt(), mean)
		diffFromMeanSquared := diffFromMean.Exp(diffFromMean, big.NewInt(2), nil)
		standardDev.Add(standardDev, diffFromMeanSquared)
	}

	return new(big.Int).Sqrt(new(big.Int).Div(standardDev, arrayLengthBigInt))
}

func getPricesArrayWithoutOutliers(prices []sdk.Int, mean *big.Int, standardDev *big.Int) []big.Int {
	var cleansedPriceArray []big.Int
	for i := 0; i < len(prices); i++ {
		ubtPrice := prices[i].BigInt()

		oneStDevLessFromMean := new(big.Int).Sub(mean, standardDev)
		oneStDevGreaterFromMean := new(big.Int).Add(mean, standardDev)

		if (ubtPrice.Cmp(oneStDevLessFromMean) == +1 || ubtPrice.Cmp(oneStDevLessFromMean) == 0) &&
			(ubtPrice.Cmp(oneStDevGreaterFromMean) == -1 || ubtPrice.Cmp(oneStDevGreaterFromMean) == 0) {
			cleansedPriceArray = append(cleansedPriceArray, *ubtPrice)
		}
	}

	return cleansedPriceArray
}

func calcAvgPrice(prices []big.Int) *big.Int {
	sum := big.NewInt(0)

	for i := 0; i < len(prices); i++ {
		sum.Add(sum, &prices[i])
	}

	arrayLengthBigInt := big.NewInt(int64(len(prices)))

	return new(big.Int).Div(sum, arrayLengthBigInt)
}
