# Blockcalc

Calculate a future block height at a given date.
Mainly used to estimate a height for network upgrades.

## Usage 

Build:

```bash
go build ./cmd/blockcalc/...
```

Run:

```bash
./blockcalc -help
Usage of ./blockcalc:
  -node string
        rpc endpoint (default "https://rpc-fetchhub.fetch.ai:443")
  -numblocks int
        number of blocks for average duration calculation (default 50)
  -target string
        the target time (format 2006-01-02T15:04:05Z07:00)
```