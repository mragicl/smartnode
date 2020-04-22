package minipool

import (
    "errors"
    "fmt"
    "strings"

    "github.com/urfave/cli"

    minipoolapi "github.com/rocket-pool/smartnode/shared/api/minipool"
    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/services/rocketpool/minipool"
    "github.com/rocket-pool/smartnode/shared/services/rocketpool/settings"
    "github.com/rocket-pool/smartnode/shared/utils/eth"
)


// Get the node's minipool statuses
func getMinipoolStatus(c *cli.Context, statusFilters []string) error {

    // Initialise services
    p, err := services.NewProvider(c, services.ProviderOpts{
        AM: true,
        CM: true,
        LoadContracts: []string{"rocketMinipoolSettings", "rocketPoolToken", "utilAddressSetStorage"},
        LoadAbis: []string{"rocketMinipool"},
        WaitClientConn: true,
        WaitClientSync: true,
        WaitRocketStorage: true,
    })
    if err != nil { return err }
    defer p.Cleanup()

    /*
    // Get latest block header
    header, err := p.Client.HeaderByNumber(context.Background(), nil)
    if err != nil {
        return errors.New("Error retrieving latest block header: " + err.Error())
    }
    */

    // Get minipool statuses
    status, err := minipoolapi.GetMinipoolStatus(p)
    if err != nil { return err }

    // Filter minipool details
    filteredMinipoolDetails := []*minipool.Details{}
    for _, details := range status.Minipools {
        statusMatch := (len(statusFilters) == 0)
        if !statusMatch {
            for _, statusFilter := range statusFilters {
                if details.StatusType == statusFilter {
                    statusMatch = true
                    break
                }
            }
        }
        if statusMatch {
            filteredMinipoolDetails = append(filteredMinipoolDetails, details)
        }
    }

    // Get minipool staking durations
    stakingDurations, err := settings.GetMinipoolStakingDurations(p.CM)
    if err != nil {
        return errors.New("Error retrieving minipool staking durations: " + err.Error())
    }

    // Get minipool overview details
    overview := map[string]map[string]uint64{"total": map[string]uint64{"total": 0}}
    for _, duration := range stakingDurations { overview["total"][duration.Id] = 0 }
    for _, details := range filteredMinipoolDetails {
        if _, ok := overview[details.StatusType]; !ok {
            overview[details.StatusType] = map[string]uint64{"total": 0}
            for _, duration := range stakingDurations { overview[details.StatusType][duration.Id] = 0 }
        }
        overview[details.StatusType][details.StakingDurationId]++
        overview[details.StatusType]["total"]++
        overview["total"][details.StakingDurationId]++
        overview["total"]["total"]++
    }

    // Print minipool overview
    if len(filteredMinipoolDetails) > 0 {

        // Header
        fmt.Fprintln(p.Output, "========================")
        fmt.Fprintln(p.Output, "Node minipools overview:")
        fmt.Fprintln(p.Output, "========================")
        fmt.Fprintln(p.Output, "")

        // Heading row
        rowText := fmt.Sprintf("%-13v | ", "Status")
        for _, duration := range stakingDurations { rowText += fmt.Sprintf("%-3v | ", duration.Id) }
        rowText += fmt.Sprintf("%v", "Total")
        fmt.Fprintln(p.Output, rowText)

        // Divider
        rowLength := 16
        for _, _ = range stakingDurations { rowLength += 6 }
        rowLength += 5
        fmt.Fprintln(p.Output, strings.Repeat("-", rowLength))

        // Content
        rowsPrinted := 0
        for _, status := range []string{"initialized", "depositassigned", "prelaunch", "staking", "loggedout", "withdrawn", "timedout", "total"} {
            if _, ok := overview[status]; !ok { continue }
            if status == "total" && rowsPrinted < 2 { continue }

            // Row
            rowText := fmt.Sprintf("%-13v | ", strings.Title(status) + ":")
            for _, duration := range stakingDurations { rowText += fmt.Sprintf("%-3v | ", overview[status][duration.Id]) }
            rowText += fmt.Sprintf("%v", overview[status]["total"])
            fmt.Fprintln(p.Output, rowText)

            rowsPrinted++
        }
        fmt.Fprintln(p.Output, "")

    }

    // Print minipool statuses
    poolsType := strings.Join(statusFilters, " / ")
    if poolsType == "" { poolsType = "total" }
    poolsTitle := fmt.Sprintf("Node has %d %s minipools:", len(filteredMinipoolDetails), poolsType)
    fmt.Fprintln(p.Output, strings.Repeat("=", len(poolsTitle)))
    fmt.Fprintln(p.Output, poolsTitle)
    fmt.Fprintln(p.Output, strings.Repeat("=", len(poolsTitle)))
    for _, details := range filteredMinipoolDetails {

        /*
        // Get staking info
        var stakingBlocksLeft int64
        var stakingCompleteAt time.Time
        if details.StakingExitBlock != nil {
            stakingBlocksLeft = details.StakingExitBlock.Int64() - header.Number.Int64()
            if stakingBlocksLeft < 0 { stakingBlocksLeft = 0 }
            stakingTimeLeft, _ := time.ParseDuration(fmt.Sprintf("%dm", stakingBlocksLeft / 4))
            stakingCompleteAt = time.Now().Add(stakingTimeLeft)
        }
        */

        // Print
        fmt.Fprintln(p.Output, "")
        fmt.Fprintln(p.Output, "Address:                 ", details.Address.Hex())
        fmt.Fprintln(p.Output, "Status:                  ", strings.Title(details.StatusType))
        fmt.Fprintln(p.Output, "Status Updated @ Time:   ", details.StatusTime.Format("2006-01-02, 15:04 -0700 MST"))
        fmt.Fprintln(p.Output, "Status Updated @ Block:  ", details.StatusBlock.String())
        fmt.Fprintln(p.Output, "")
        fmt.Fprintln(p.Output, "Staking Duration:        ", details.StakingDurationId)
        fmt.Fprintln(p.Output, "Staking Total Epochs:    ", details.StakingDuration.String())

        /*
        if details.StakingExitBlock != nil {
        fmt.Fprintln(p.Output, "Staking Until Block:     ", details.StakingExitBlock.String())
        fmt.Fprintln(p.Output, "Staking Blocks Left:     ", stakingBlocksLeft)
        fmt.Fprintln(p.Output, "Staking Complete Approx: ", stakingCompleteAt.Format("2006-01-02, 15:04 -0700 MST"))
        }
        */

        fmt.Fprintln(p.Output, "")
        if details.Status >= minipool.WITHDRAWN {
        fmt.Fprintln(p.Output, "Node Deposit Withdrawn:  ", fmt.Sprintf("%t", !details.NodeDepositExists))
        }
        fmt.Fprintln(p.Output, "Node ETH Deposited:      ", fmt.Sprintf("%.2f", eth.WeiToEth(details.NodeEtherBalanceWei)))
        fmt.Fprintln(p.Output, "Node RPL Deposited:      ", fmt.Sprintf("%.2f", eth.WeiToEth(details.NodeRplBalanceWei)))
        fmt.Fprintln(p.Output, "")
        fmt.Fprintln(p.Output, "User Deposit Count:      ", details.UserDepositCount.String())
        if details.Status <= minipool.STAKING || details.Status == minipool.TIMED_OUT {
        fmt.Fprintln(p.Output, "User Deposit Total:      ", fmt.Sprintf("%.2f", eth.WeiToEth(details.UserDepositTotalWei)))
        fmt.Fprintln(p.Output, "User Deposit Capacity:   ", fmt.Sprintf("%.2f", eth.WeiToEth(details.UserDepositCapacityWei)))
        }
        fmt.Fprintln(p.Output, "")
        fmt.Fprintln(p.Output, "--------")

    }

    // Return
    return nil

}
