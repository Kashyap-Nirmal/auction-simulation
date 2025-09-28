# Auction Simulator

A concurrent auction simulator written in Go.  
Simulates multiple auctions with bidders responding to object attributes, enforces timeouts, declares winners, and measures total execution time.

---

## Features
- 100 simulated bidders.
- Each auction is defined by 20 attributes.
- Bidders may or may not respond with a bid.
- Auctions run with timeouts — only bids before timeout are considered.
- 40 auctions run concurrently.
- Reports execution time (start of first auction → end of last auction).
- Simple resource standardization via environment variables (`MAX_VCPU`, `MAX_RAM`).

---

```
## Project Structure

auction-simulator/
├── cmd/
│   └── main.go        # Entry point
├── auction/
│   └── auction.go     # Auction logic, bids, results
├── bidder/
│   └── bidder.go      # Bidder simulation
├── util/
│   └── env.go         # Resource configuration
└── go.mod
```
---

## Usage

### Prerequisites
- Go 1.20+ installed.

### Run
```bash
go run ./cmd/main.go
````

### Configure Resources

```bash
export SIM_VCPU=2
export MAX_RAM=1024
go run ./cmd/main.go
```

---

# auction-simulation
