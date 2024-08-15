package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetTotalSwap(address string, startTime time.Time, endTime time.Time) (float64, error) {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		return 0, err
	}
	defer client.Close()

	uniswapV2Address := common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{uniswapV2Address},
		Topics: [][]common.Hash{
			{common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822")},
		},
		FromBlock: big.NewInt(12345678),
		ToBlock:   big.NewInt(987654321),
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case log := <-logs:
			fmt.Println(log)
		}
	}

	return 0, nil
}
