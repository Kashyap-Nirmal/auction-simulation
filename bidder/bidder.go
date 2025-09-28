package bidder

import (
	"math/rand"
	"time"

	"auction-simulator/auction"
)

// NewBidder returns a closure that simulates a bidder.
// It matches the type used by auction.RunAuction: func([]string) *auction.Bid
func NewBidder(id int) func([]string) *auction.Bid {
	return func(attrs []string) *auction.Bid {
		// Simulate processing/network delay (up to 150ms)
		time.Sleep(time.Duration(rand.Intn(150)) * time.Millisecond)

		// 30% chance the bidder does not respond
		if rand.Float64() < 0.3 {
			return nil
		}

		// Create a random bid amount (tune as needed)
		amount := 50 + rand.Float64()*100
		return &auction.Bid{
			BidderID: id,
			Amount:   amount,
		}
	}
}
