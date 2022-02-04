package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Block struct {
	Height int64
	Time   time.Time
}

func getBlock(node string, height int64) (*Block, error) {
	u, err := url.Parse(fmt.Sprintf("%s/block", node))
	if err != nil {
		return nil, err
	}

	if height > 0 {
		v := url.Values{}
		v.Add("height", fmt.Sprintf("%d", height))
		u.RawQuery = v.Encode()
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	blockResp := &struct {
		Result struct {
			Block struct {
				Header struct {
					Height string    `json:"height"`
					Time   time.Time `json:"time"`
				} `json:"header"`
			} `json:"block"`
		} `json:"result"`
	}{}

	if err := dec.Decode(blockResp); err != nil {
		return nil, err
	}

	blockHeight, err := strconv.ParseInt(blockResp.Result.Block.Header.Height, 10, 64)
	if err != nil {
		return nil, err
	}

	return &Block{
		Height: blockHeight,
		Time:   blockResp.Result.Block.Header.Time,
	}, nil
}

func getCurrentBlock(node string) (*Block, error) {
	return getBlock(node, 0)
}

func getAverageBlockTime(node string, numBlocks int) (time.Duration, error) {
	current, err := getCurrentBlock(node)
	if err != nil {
		return 0, err
	}

	startHeight := current.Height

	var total int64
	for i := 1; i <= numBlocks; i++ {
		h := startHeight - int64(i)
		previous, err := getBlock(node, h)
		if err != nil {
			return 0, err
		}
		total += current.Time.Sub(previous.Time).Nanoseconds()
		current = previous
	}

	return time.Duration(total / int64(numBlocks)), nil
}

func main() {
	var targetTime string
	var rpcEndpoint string
	var numBlocks int

	flag.StringVar(&targetTime, "target", "", fmt.Sprintf("the target time (format %s)", time.RFC3339))
	flag.StringVar(&rpcEndpoint, "node", "https://rpc-fetchhub.fetch.ai:443", "rpc endpoint")
	flag.IntVar(&numBlocks, "numblocks", 50, "number of blocks for average duration calculation")
	flag.Parse()

	target, err := time.Parse(time.RFC3339, targetTime)
	if err != nil {
		panic(err)
	}

	currentBlock, err := getCurrentBlock(rpcEndpoint)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current block height: %d\n", currentBlock.Height)

	avg, err := getAverageBlockTime(rpcEndpoint, numBlocks)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Average block time: %s\n", avg)

	diff := target.Sub(currentBlock.Time)

	nbBlocks := diff.Seconds() / avg.Seconds()
	fmt.Printf("Target is in %d blocks (%s)\n", int64(nbBlocks), diff)
	fmt.Printf("Estimated block height at %s: %d\n", target, currentBlock.Height+int64(nbBlocks))
}
