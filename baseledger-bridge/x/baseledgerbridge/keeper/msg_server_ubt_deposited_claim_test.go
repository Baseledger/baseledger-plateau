package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Ognjen - Add the real method for avg calculation here
func TestAveragePriceCalculation(t *testing.T) {
	pricesArray := []string{}
	existingPrice := "0.80"
	newPrice := "0.70"

	existingPriceDec, _ := sdk.NewDecFromStr(existingPrice)
	newPriceDec, _ := sdk.NewDecFromStr(newPrice)

	existingPriceInt := sdk.NewIntFromBigInt(existingPriceDec.BigInt())
	newPriceInt := sdk.NewIntFromBigInt(newPriceDec.BigInt())

	existingAvgPrice := existingPriceInt

	pricesArray = append(pricesArray, existingPriceInt.String())
	pricesArray = append(pricesArray, newPriceInt.String())

	ubtPricesNewLength := sdk.NewInt(int64((len(pricesArray))))
	avgAddition := newPriceInt.Sub(existingAvgPrice).Quo(ubtPricesNewLength)
	newAvgPrice := existingAvgPrice.Add(avgAddition)

	require.Equal(t, sdk.NewInt(750000000000000000), newAvgPrice)
}
