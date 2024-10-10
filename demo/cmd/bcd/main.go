package main

import (
	"os"

	"github.com/babylonlabs-io/babylon-sdk/demo/app"
	"github.com/cometbft/cometbft/libs/log"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		log.NewTMLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
