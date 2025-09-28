package auction

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Bid struct {
	BidderID int
	Amount   float64
	Time     time.Time
}

type BidResult struct {
	AuctionID int
	Bids      []Bid
	WinnerID  int
	Amount    float64
}

func (br BidResult) WinnerInfo() string {
	if br.WinnerID == -1 {
		return "No valid bids"
	}
	return fmt.Sprintf("Bidder %2d with amount %.2f", br.WinnerID, br.Amount)
}

type Auction struct {
	ID         int
	Attributes []string
	Timeout    time.Duration
}

// GenerateAuction creates an auction with n attributes and a timeout
func GenerateAuction(id int, numAttrs int, timeout time.Duration) Auction {
	attrs := make([]string, numAttrs)
	for i := range attrs {
		attrs[i] = fmt.Sprintf("Attr-%d", i+1)
	}
	return Auction{ID: id, Attributes: attrs, Timeout: timeout}
}

// RunAuction runs the auction with given bidders and returns BidResult
func RunAuction(ctx context.Context, a Auction, bidders []func([]string) *Bid) BidResult {
    bidCh := make(chan *Bid, len(bidders))
    var wg sync.WaitGroup
    var allBids []Bid
    var bidsMutex sync.Mutex // Add this line

    for _, bidder := range bidders {
        wg.Add(1)
        go func(b func([]string) *Bid) {
            defer wg.Done()
            if bid := b(a.Attributes); bid != nil {
                bidsMutex.Lock()                    // Add this line
                allBids = append(allBids, *bid)     // Now thread-safe
                bidsMutex.Unlock()                  // Add this line
                
                select {
                case bidCh <- bid:
                case <-ctx.Done():
                }
            }
        }(bidder)
    }

    // Rest of your code stays exactly the same...
    go func() {
        wg.Wait()
        close(bidCh)
    }()

    timer := time.NewTimer(a.Timeout)
    defer timer.Stop()

    var highest *Bid
    for {
        select {
        case bid, ok := <-bidCh:
            if !ok {
                return finalize(a.ID, highest, allBids)
            }
            if highest == nil || bid.Amount > highest.Amount {
                highest = bid
            }
        case <-timer.C:
            return finalize(a.ID, highest, allBids)
        }
    }
}

// finalize returns BidResult with all bids
func finalize(auctionID int, highest *Bid, allBids []Bid) BidResult {
	if highest == nil {
		return BidResult{AuctionID: auctionID, WinnerID: -1, Amount: 0, Bids: allBids}
	}
	return BidResult{
		AuctionID: auctionID,
		WinnerID:  highest.BidderID,
		Amount:    highest.Amount,
		Bids:      allBids,
	}
}
