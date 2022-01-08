package migrations

import (
	"log"

	"github.com/Baseledger/baseledger-bridge/x/baseledgerbridge/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper keeper.Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper keeper.Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	log.Print("Migrate2to3 WORKS")
	return nil
}
