package app

import (
	"fmt"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/stretchr/testify/require"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var emptyWasmOpts []wasm.Option

// adapted from https://github.com/cosmos/cosmos-sdk/blob/v0.50.6/simapp/app_test.go#L47-L48
func TestSimAppExportAndBlockedAddrs(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()
	app := NewAppWithCustomOptions(t, false, SetupOptions{
		Logger:  logger.With("instance", "first"),
		DB:      db,
		AppOpts: simtestutil.NewAppOptionsWithFlagHome(t.TempDir()),
	})

	// BlockedAddresses returns a map of addresses in app v1 and a map of modules name in app v2.
	for acc := range BlockedAddresses() {
		var addr sdk.AccAddress
		if modAddr, err := sdk.AccAddressFromBech32(acc); err == nil {
			addr = modAddr
		} else {
			addr = app.AccountKeeper.GetModuleAddress(acc)
		}

		require.True(
			t,
			app.BankKeeper.BlockedAddr(addr),
			fmt.Sprintf("ensure that blocked addresses are properly set in bank keeper: %s should be blocked", acc),
		)
	}

	// finalize block so we have CheckTx state set
	//_, err := app.FinalizeBlock(&abci.RequestFinalizeBlock{
	//	Height: 1,
	//})

	_ = app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewConsumerApp(logger, db, nil, true, simtestutil.NewAppOptionsWithFlagHome(t.TempDir()), emptyWasmOpts)
	_, err := app2.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}
