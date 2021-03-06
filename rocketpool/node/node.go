package node

import (
    "github.com/fatih/color"
    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/utils/log"
)


// Config
const (
    StakePrelaunchMinipoolsColor = color.FgBlue
)


// Register node command
func RegisterCommands(app *cli.App, name string, aliases []string) {
    app.Commands = append(app.Commands, cli.Command{
        Name:      name,
        Aliases:   aliases,
        Usage:     "Run Rocket Pool node activity daemon",
        Action: func(c *cli.Context) error {
            return run(c)
        },
    })
}


// Run daemon
func run(c *cli.Context) error {

    // Wait until node is registered
    if err := services.WaitNodeRegistered(c, true); err != nil { return err }

    // Initialize tasks
    stakePrelaunchMinipools, err := newStakePrelaunchMinipools(c, log.NewColorLogger(StakePrelaunchMinipoolsColor))
    if err != nil { return err }

    // Start tasks
    stakePrelaunchMinipools.Start()

    // Block thread
    select {}

}

