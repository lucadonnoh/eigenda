package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/subgraph"
	"github.com/Layr-Labs/eigenda/tools/opscan"
	"github.com/Layr-Labs/eigenda/tools/opscan/flags"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	version   = ""
	gitCommit = ""
	gitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s,%s,%s", version, gitCommit, gitDate)
	app.Name = "opscan"
	app.Description = "operator network scanner"
	app.Usage = ""
	app.Flags = flags.Flags
	app.Action = RunScan
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RunScan(ctx *cli.Context) error {
	config, err := opscan.NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(config.LoggerConfig)
	if err != nil {
		return err
	}

	subgraphApi := subgraph.NewApi(config.SubgraphEndpoint, config.SubgraphEndpoint)
	subgraphClient := dataapi.NewSubgraphClient(subgraphApi, logger)

	semvers := make(map[string]int)
	if config.OperatorId != "" {
		operatorInfo, err := subgraphClient.QueryOperatorInfoByOperatorId(context.Background(), config.OperatorId)
		if err != nil {
			logger.Warn("failed to fetch operator info", "operatorId", config.OperatorId, "error", err)
			return errors.New("operator info not found")
		}

		operatorSocket := core.OperatorSocket(operatorInfo.Socket)
		retrievalSocket := operatorSocket.GetRetrievalSocket()
		semver := getNodeInfo(context.Background(), retrievalSocket, config.Timeout, logger)
		semvers[semver]++

	} else {
		indexedOperatorState, err := subgraphClient.QueryOperatorsWithLimit(context.Background(), 1000)
		if err != nil {
			return fmt.Errorf("failed to fetch indexed operator state - %s", err)
		}
		logger.Info("Scanning operators", "count", len(indexedOperatorState))
		//semvers = scanOperators(indexedOperatorState, config, logger)
	}
	displayResults(semvers)
	return nil
}

func scanOperators(indexedOperatorState *dataapi.IndexedQueriedOperatorInfo, config *opscan.Config, logger logging.Logger) map[string]int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	semvers := make(map[string]int)
	for _, operator := range indexedOperatorState.Operators {
		wg.Add(1)
		go func(operator dataapi.QueriedOperatorInfo) {
			defer wg.Done()
			operatorSocket := core.OperatorSocket(operator.IndexedOperatorInfo.Socket)
			retrievalSocket := operatorSocket.GetRetrievalSocket()
			semver := getNodeInfo(context.Background(), retrievalSocket, config.Timeout, logger)

			mu.Lock()
			semvers[semver]++
			mu.Unlock()
		}(*operator)
	}

	wg.Wait()
	return semvers
}

func getNodeInfo(ctx context.Context, socket string, timeout time.Duration, logger logging.Logger) string {
	return "dry-run"
	conn, err := grpc.Dial(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to dial grpc operator socket", "socket", socket, "error", err)
		return "unreachable"
	}
	defer conn.Close()
	client := node.NewRetrievalClient(conn)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	reply, err := client.NodeInfo(ctx, &node.NodeInfoRequest{})
	if err != nil {
		var semver string
		if strings.Contains(err.Error(), "Unimplemented") {
			semver = "<0.8.0"
		} else if strings.Contains(err.Error(), "DeadlineExceeded") {
			semver = "timeout"
		} else if strings.Contains(err.Error(), "Unavailable") {
			semver = "refused"
		}
		logger.Warn("NodeInfo", "semver", semver, "error", err)
		return semver
	}

	logger.Info("NodeInfo", "semver", reply.Semver, "os", reply.Os, "arch", reply.Arch, "numCpu", reply.NumCpu, "memBytes", reply.MemBytes)
	return reply.Semver
}

func displayResults(results map[string]int) {
	tw := table.NewWriter()

	rowHeader := table.Row{"semver", "count"}
	tw.AppendHeader(rowHeader)

	total := 0
	for semver, count := range results {
		tw.AppendRow(table.Row{semver, count})
		total += count
	}
	tw.AppendFooter(table.Row{"total", total})

	fmt.Println(tw.Render())
}