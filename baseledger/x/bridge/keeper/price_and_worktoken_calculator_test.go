package keeper

import (
	"math/big"
	"testing"

	"github.com/Baseledger/baseledger/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCalculateAvgUbtPriceForAttestation_TwoPriceValues_ReturnsAverage(t *testing.T) {
	// Arrange
	var attestation = types.Attestation{
		Observed:  false,
		Votes:     []string{},
		Height:    uint64(1),
		UbtPrices: []sdk.Int{},
	}

	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(800000000000000000))
	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(700000000000000000))

	// Act
	avgPrice := CalculateAvgUbtPriceForAttestation(attestation)

	// Assert
	require.Equal(t, big.NewInt(750000000000000000), avgPrice)
}

func TestCalculateAvgUbtPriceForAttestation_ThreePriceValues__ReturnsAverage(t *testing.T) {
	// Arrange
	var attestation = types.Attestation{
		Observed:  false,
		Votes:     []string{},
		Height:    uint64(1),
		UbtPrices: []sdk.Int{},
	}

	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(800000000000000000))
	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(700000000000000000))
	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(600000000000000000))

	// Act
	avgPrice := CalculateAvgUbtPriceForAttestation(attestation)

	// Assert
	require.Equal(t, big.NewInt(700000000000000000), avgPrice)
}

func TestCalculateAvgUbtPriceForAttestation_FourPriceValuesOneOutlier__ReturnsAverageWithoutOutlier(t *testing.T) {
	// Arrange
	var attestation = types.Attestation{
		Observed:  false,
		Votes:     []string{},
		Height:    uint64(1),
		UbtPrices: []sdk.Int{},
	}

	outlier, _ := sdk.NewIntFromString("99900000000000000000")

	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(800000000000000000))
	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(700000000000000000))
	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(600000000000000000))
	attestation.UbtPrices = append(attestation.UbtPrices, outlier)

	// Act
	avgPrice := CalculateAvgUbtPriceForAttestation(attestation)

	// Assert
	require.Equal(t, big.NewInt(700000000000000000), avgPrice)
}

func TestCalculateAmountOfWorkTokens_UbtAmountAndPriceWithoutMod__ReturnsCorrectAmount(t *testing.T) {
	// Arrange
	depositedUbtAmount, _ := new(big.Int).SetString("10000000000", 0) // 100
	avgUbtPrice, _ := new(big.Int).SetString("100000000", 0)          // 1
	worktokenEurPrice, _ := sdk.NewDecFromStr("0.1")

	// Act
	calculatedWorkTokenAmount := CalculateAmountOfWorkTokens(worktokenEurPrice.BigInt(), depositedUbtAmount, avgUbtPrice)

	// Assert
	expectedWorktokenAmount, _ := new(big.Int).SetString("1000", 0)
	require.Equal(t, expectedWorktokenAmount, calculatedWorkTokenAmount)
}

func TestCalculateAmountOfWorkTokens_UbtAmountAndPriceWithMod__ReturnsCeiledAmount(t *testing.T) {
	// Arrange
	depositedUbtAmount, _ := new(big.Int).SetString("10001000000", 0) // 100.01
	avgUbtPrice, _ := new(big.Int).SetString("100000000", 0)          // 1
	worktokenEurPrice, _ := sdk.NewDecFromStr("0.1")

	// Act
	calculatedWorkTokenAmount := CalculateAmountOfWorkTokens(worktokenEurPrice.BigInt(), depositedUbtAmount, avgUbtPrice)

	// Assert
	expectedWorktokenAmount, _ := new(big.Int).SetString("1001", 0)
	require.Equal(t, expectedWorktokenAmount, calculatedWorkTokenAmount)
}

func TestCalculateAmountOfWorkTokens_UbtAmountAndPriceWithModVariation__ReturnsCeiledAmount(t *testing.T) {
	// Arrange
	depositedUbtAmount, _ := new(big.Int).SetString("9999000000", 0) // 99.99
	avgUbtPrice, _ := new(big.Int).SetString("100000000", 0)         // 1
	worktokenEurPrice, _ := sdk.NewDecFromStr("0.1")

	// Act
	calculatedWorkTokenAmount := CalculateAmountOfWorkTokens(worktokenEurPrice.BigInt(), depositedUbtAmount, avgUbtPrice)

	// Assert
	expectedWorktokenAmount, _ := new(big.Int).SetString("1000", 0)
	require.Equal(t, expectedWorktokenAmount, calculatedWorkTokenAmount)
}

func TestCalculateAmountOfWorkTokens_UbtAmountJustAboveOneTokenAndUbtPriceOneEur__ReturnsCeiledAmount(t *testing.T) {
	// Arrange
	depositedUbtAmount, _ := new(big.Int).SetString("100000001", 0) // 1.00000001
	avgUbtPrice, _ := new(big.Int).SetString("100000000", 0)         // 1
	worktokenEurPrice, _ := sdk.NewDecFromStr("0.1")

	// Act
	calculatedWorkTokenAmount := CalculateAmountOfWorkTokens(worktokenEurPrice.BigInt(), depositedUbtAmount, avgUbtPrice)

	// Assert
	expectedWorktokenAmount, _ := new(big.Int).SetString("11", 0)
	require.Equal(t, expectedWorktokenAmount, calculatedWorkTokenAmount)
}

func TestCalculateAmountOfWorkTokens_UbtAmountJustBellowTwoTokensAndUbtPriceOneEur__ReturnsCeiledAmount(t *testing.T) {
	// Arrange
	depositedUbtAmount, _ := new(big.Int).SetString("199999999", 0) // 1.99999999
	avgUbtPrice, _ := new(big.Int).SetString("100000000", 0)         // 1
	worktokenEurPrice, _ := sdk.NewDecFromStr("0.1")

	// Act
	calculatedWorkTokenAmount := CalculateAmountOfWorkTokens(worktokenEurPrice.BigInt(), depositedUbtAmount, avgUbtPrice)

	// Assert
	expectedWorktokenAmount, _ := new(big.Int).SetString("20", 0)
	require.Equal(t, expectedWorktokenAmount, calculatedWorkTokenAmount)
}

func TestCalculateAmountOfWorkTokens_UbtAmountJustBellowTwoTokensAndUbtPriceZeroNineEur__ReturnsCeiledAmount(t *testing.T) {
	// Arrange
	depositedUbtAmount, _ := new(big.Int).SetString("199999999", 0) // 1.99999999
	avgUbtPrice, _ := new(big.Int).SetString("90000000", 0)         // 0.9
	worktokenEurPrice, _ := sdk.NewDecFromStr("0.1")

	// Act
	calculatedWorkTokenAmount := CalculateAmountOfWorkTokens(worktokenEurPrice.BigInt(), depositedUbtAmount, avgUbtPrice)

	// Assert
	expectedWorktokenAmount, _ := new(big.Int).SetString("18", 0)
	require.Equal(t, expectedWorktokenAmount, calculatedWorkTokenAmount)
}
