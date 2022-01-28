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
	depositedUbtAmount, _ := new(big.Int).SetString("100000000000000000000", 0) // 100
	avgUbtPrice, _ := new(big.Int).SetString("1000000000000000000", 0)          // 1

	// Act
	calculatedWorkTokenAmount := CalculateAmountOfWorkTokens(depositedUbtAmount, avgUbtPrice)

	// Assert
	expectedWorktokenAmount, _ := new(big.Int).SetString("1000000000000000000000", 0) // 1000
	require.Equal(t, expectedWorktokenAmount, calculatedWorkTokenAmount)
}

func TestCalculateAmountOfWorkTokens_UbtAmountAndPriceWithMod__ReturnsCeiledAmount(t *testing.T) {
	// Arrange
	depositedUbtAmount, _ := new(big.Int).SetString("100010000000000000000", 0) // 100.01
	avgUbtPrice, _ := new(big.Int).SetString("1000000000000000000", 0)          // 1

	// Act
	calculatedWorkTokenAmount := CalculateAmountOfWorkTokens(depositedUbtAmount, avgUbtPrice)

	// Assert
	expectedWorktokenAmount, _ := new(big.Int).SetString("1001000000000000000000", 0) // 1001
	require.Equal(t, expectedWorktokenAmount, calculatedWorkTokenAmount)
}

func TestCalculateAmountOfWorkTokens_UbtAmountAndPriceWithModVariation__ReturnsCeiledAmount(t *testing.T) {
	// Arrange
	depositedUbtAmount, _ := new(big.Int).SetString("99990000000000000000", 0) // 99.99
	avgUbtPrice, _ := new(big.Int).SetString("1000000000000000000", 0)         // 1

	// Act
	calculatedWorkTokenAmount := CalculateAmountOfWorkTokens(depositedUbtAmount, avgUbtPrice)

	// Assert
	expectedWorktokenAmount, _ := new(big.Int).SetString("1000000000000000000000", 0) // 1000
	require.Equal(t, expectedWorktokenAmount, calculatedWorkTokenAmount)
}
