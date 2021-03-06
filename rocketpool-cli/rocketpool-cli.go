package main

import (
    "fmt"
    "os"

    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/rocketpool-cli/faucet"
    "github.com/rocket-pool/smartnode/rocketpool-cli/minipool"
    "github.com/rocket-pool/smartnode/rocketpool-cli/network"
    "github.com/rocket-pool/smartnode/rocketpool-cli/node"
    "github.com/rocket-pool/smartnode/rocketpool-cli/queue"
    "github.com/rocket-pool/smartnode/rocketpool-cli/service"
    "github.com/rocket-pool/smartnode/rocketpool-cli/wallet"
)


// Run
func main() {

    // Add logo to application help template
    cli.AppHelpTemplate = fmt.Sprintf(`
______           _        _    ______           _ 
| ___ \         | |      | |   | ___ \         | |
| |_/ /___   ___| | _____| |_  | |_/ /__   ___ | |
|    // _ \ / __| |/ / _ \ __| |  __/ _ \ / _ \| |
| |\ \ (_) | (__|   <  __/ |_  | | | (_) | (_) | |
\_| \_\___/ \___|_|\_\___|\__| \_|  \___/ \___/|_|

%s`, cli.AppHelpTemplate)

    // Initialise application
    app := cli.NewApp()

    // Set application info
    app.Name = "rocketpool"
    app.Usage = "Rocket Pool CLI"
    app.Version = "0.0.1"
    app.Authors = []cli.Author{
        cli.Author{
            Name:  "David Rugendyke",
            Email: "david@rocketpool.net",
        },
        cli.Author{
            Name:  "Jake Pospischil",
            Email: "jake@rocketpool.net",
        },
    }
    app.Copyright = "(c) 2020 Rocket Pool Pty Ltd"

    // Set application flags
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "host, o",
            Usage: "Smart node SSH host `address`",
        },
        cli.StringFlag{
            Name:  "user, u",
            Usage: "Smart node SSH user `name`",
        },
        cli.StringFlag{
            Name:  "key, k",
            Usage: "Smart node SSH key `file`",
        },
    }

    // Register commands
      faucet.RegisterCommands(app, "faucet",   []string{"f"})
    minipool.RegisterCommands(app, "minipool", []string{"m"})
     network.RegisterCommands(app, "network",  []string{"e"})
        node.RegisterCommands(app, "node",     []string{"n"})
       queue.RegisterCommands(app, "queue",    []string{"q"})
     service.RegisterCommands(app, "service",  []string{"s"})
      wallet.RegisterCommands(app, "wallet",   []string{"w"})

    // Run application
    fmt.Println("")
    if err := app.Run(os.Args); err != nil {
        fmt.Println(err)
    }
    fmt.Println("")

}

