package keeper

import (
	"math/big"
	"testing"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestTwoPriceValues_CalculateAvgUbtPriceForAttestation_ReturnsAverage(t *testing.T) {
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

func TestThreePriceValues_CalculateAvgUbtPriceForAttestation_ReturnsAverage(t *testing.T) {
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

func TestFourPriceValuesOneOutlier_CalculateAvgUbtPriceForAttestation_ReturnsAverageWithoutOutlier(t *testing.T) {
	// Arrange
	var attestation = types.Attestation{
		Observed:  false,
		Votes:     []string{},
		Height:    uint64(1),
		UbtPrices: []sdk.Int{},
	}

	outlier, _ := sdk.NewIntFromString("99600000000000000000")

	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(800000000000000000))
	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(700000000000000000))
	attestation.UbtPrices = append(attestation.UbtPrices, sdk.NewInt(600000000000000000))
	attestation.UbtPrices = append(attestation.UbtPrices, outlier)

	// Act
	avgPrice := CalculateAvgUbtPriceForAttestation(attestation)

	// Assert
	require.Equal(t, big.NewInt(700000000000000000), avgPrice)
}
