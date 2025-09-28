// package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"math/rand"
// 	"os"
// 	"runtime"
// 	"sort"
// 	"sync"
// 	"time"

// 	"auction-simulator/auction"
// 	"auction-simulator/bidder"
// 	"auction-simulator/util"
// )

// func main() {
// 	rand.Seed(time.Now().UnixNano())

// 	vcpu := util.GetEnvInt("SIM_VCPU", runtime.NumCPU())
// 	runtime.GOMAXPROCS(vcpu)
// 	fmt.Printf("Using GOMAXPROCS=%d\n\n", vcpu)

// 	numAuctions := 40
// 	numBidders := 100

// 	// Prepare bidders
// 	bidders := make([]func([]string) *auction.Bid, numBidders)
// 	for i := 0; i < numBidders; i++ {
// 		bidders[i] = bidder.NewBidder(i + 1)
// 	}

// 	var wg sync.WaitGroup
// 	resultsCh := make(chan auction.BidResult, numAuctions)
// 	var startFirstTime time.Time
// 	var once sync.Once

// 	// Run auctions concurrently
// 	for a := 0; a < numAuctions; a++ {
// 		au := auction.GenerateAuction(a+1, 20, 1500*time.Millisecond)
// 		wg.Add(1)
// 		go func(au auction.Auction) {
// 			defer wg.Done()
// 			once.Do(func() { startFirstTime = time.Now() })
// 			res := auction.RunAuction(context.Background(), au, bidders)
// 			resultsCh <- res
// 		}(au)
// 	}

// 	wg.Wait()
// 	close(resultsCh)

// 	endLastTime := time.Now()
// 	totalTime := endLastTime.Sub(startFirstTime)
// 	fmt.Printf("Total elapsed: %v\n\n", totalTime)

// 	// Collect results
// 	var results []auction.BidResult
// 	for r := range resultsCh {
// 		results = append(results, r)
// 	}

// 	// Console output: Table sorted by highest bid
// 	fmt.Println("=== Auctions Sorted by Highest Bid ===")
// 	sort.SliceStable(results, func(i, j int) bool { return results[i].Amount > results[j].Amount })
// 	printTable(results)

// 	// Console output: Table sorted by auction number
// 	fmt.Println("\n=== Auctions Sorted by Auction ID ===")
// 	sort.SliceStable(results, func(i, j int) bool { return results[i].AuctionID < results[j].AuctionID })
// 	printTable(results)

// 	// File output: Save results to file
// 	writeResultsToFile(results, totalTime)
// }

// // Print table to console
// func printTable(results []auction.BidResult) {
// 	printTableToWriter(os.Stdout, results)
// }

// // Generic table writer - can write to console OR file
// func printTableToWriter(w io.Writer, results []auction.BidResult) {
// 	fmt.Fprintf(w, "%-8s | %-20s | %-12s | %-8s\n", "Auction", "Winner", "Total Bids", "Top Bid")
// 	fmt.Fprintf(w, "-----------------------------------------------\n")
// 	for _, r := range results {
// 		winnerStr := fmt.Sprintf("Bidder %2d", r.WinnerID)
// 		if r.WinnerID == -1 {
// 			winnerStr = "No valid bids"
// 		}
// 		fmt.Fprintf(w, "%-8d | %-20s | %-12d | %-8.2f\n", r.AuctionID, winnerStr, len(r.Bids), r.Amount)
// 	}
// }

// // Write results to file
// func writeResultsToFile(results []auction.BidResult, totalTime time.Duration) {
// 	filename := "auction_results.txt"
	
// 	// Debug: show where file will be created
// 	wd, _ := os.Getwd()
// 	fmt.Printf("\nCreating file at: %s/%s\n", wd, filename)
	
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		fmt.Printf("Error creating output file: %v\n", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Write header info
// 	fmt.Fprintf(file, "=== Auction Simulator Results ===\n")
// 	fmt.Fprintf(file, "Total Execution Time: %v\n", totalTime)
// 	fmt.Fprintf(file, "Number of Auctions: %d\n", len(results))
// 	fmt.Fprintf(file, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

// 	// Write both tables to file
// 	fmt.Fprintf(file, "=== Auctions Sorted by Highest Bid ===\n")
// 	sortedByBid := make([]auction.BidResult, len(results))
// 	copy(sortedByBid, results)
// 	sort.SliceStable(sortedByBid, func(i, j int) bool { return sortedByBid[i].Amount > sortedByBid[j].Amount })
// 	printTableToWriter(file, sortedByBid)

// 	fmt.Fprintf(file, "\n=== Auctions Sorted by Auction ID ===\n")
// 	sortedByID := make([]auction.BidResult, len(results))
// 	copy(sortedByID, results)
// 	sort.SliceStable(sortedByID, func(i, j int) bool { return sortedByID[i].AuctionID < sortedByID[j].AuctionID })
// 	printTableToWriter(file, sortedByID)

// 	fmt.Printf("\n✓ Results saved to auction_results.txt\n")
// }

package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"

	"auction-simulator/auction"
	"auction-simulator/bidder"
	"auction-simulator/util"
)

func main() {
	// Show working directory
	wd, _ := os.Getwd()
	fmt.Printf("Working directory: %s\n", wd)
	
	// List all .txt files
	files, _ := filepath.Glob("*.txt")
	fmt.Printf("All .txt files in directory: %v\n", files)
	
	fmt.Println("\n=== NOW RUNNING NORMAL AUCTION CODE ===")
	
	rand.Seed(time.Now().UnixNano())

	vcpu := util.GetEnvInt("SIM_VCPU", runtime.NumCPU())
	runtime.GOMAXPROCS(vcpu)
	fmt.Printf("Using GOMAXPROCS=%d\n\n", vcpu)

	numAuctions := 40
	numBidders := 100

	// Prepare bidders
	bidders := make([]func([]string) *auction.Bid, numBidders)
	for i := 0; i < numBidders; i++ {
		bidders[i] = bidder.NewBidder(i + 1)
	}

	var wg sync.WaitGroup
	resultsCh := make(chan auction.BidResult, numAuctions)
	var startFirstTime time.Time
	var once sync.Once

	// Run auctions concurrently
	for a := 0; a < numAuctions; a++ {
		au := auction.GenerateAuction(a+1, 20, 1500*time.Millisecond)
		wg.Add(1)
		go func(au auction.Auction) {
			defer wg.Done()
			once.Do(func() { startFirstTime = time.Now() })
			res := auction.RunAuction(context.Background(), au, bidders)
			resultsCh <- res
		}(au)
	}

	wg.Wait()
	close(resultsCh)

	endLastTime := time.Now()
	totalTime := endLastTime.Sub(startFirstTime)
	fmt.Printf("Total elapsed: %v\n\n", totalTime)

	// Collect results
	var results []auction.BidResult
	for r := range resultsCh {
		results = append(results, r)
	}

	// Console output: Table sorted by highest bid
	fmt.Println("=== Auctions Sorted by Highest Bid ===")
	sort.SliceStable(results, func(i, j int) bool { return results[i].Amount > results[j].Amount })
	printTable(results)

	// Console output: Table sorted by auction number
	fmt.Println("\n=== Auctions Sorted by Auction ID ===")
	sort.SliceStable(results, func(i, j int) bool { return results[i].AuctionID < results[j].AuctionID })
	printTable(results)

	// File output: Save results to file
	fmt.Println("=== CALLING writeResultsToFile ===")
	writeResultsToFile(results, totalTime)
	fmt.Println("=== FINISHED writeResultsToFile ===")
}

// Print table to console
func printTable(results []auction.BidResult) {
	printTableToWriter(os.Stdout, results)
}

// Generic table writer - can write to console OR file
func printTableToWriter(w io.Writer, results []auction.BidResult) {
	fmt.Fprintf(w, "%-8s | %-20s | %-12s | %-8s\n", "Auction", "Winner", "Total Bids", "Top Bid")
	fmt.Fprintf(w, "-----------------------------------------------\n")
	for _, r := range results {
		winnerStr := fmt.Sprintf("Bidder %2d", r.WinnerID)
		if r.WinnerID == -1 {
			winnerStr = "No valid bids"
		}
		fmt.Fprintf(w, "%-8d | %-20s | %-12d | %-8.2f\n", r.AuctionID, winnerStr, len(r.Bids), r.Amount)
	}
}

func writeResultsToFile(results []auction.BidResult, totalTime time.Duration) {
	fmt.Printf("\n=== Creating Individual Auction Files ===\n")
	fmt.Printf("Number of results to process: %d\n", len(results))
	
	successCount := 0
	
	// Create individual file for each auction
	for i, result := range results {
		filename := fmt.Sprintf("auction_%d_result.txt", result.AuctionID)
		fmt.Printf("Creating file %d/%d: %s\n", i+1, len(results), filename)
		
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Failed to create %s: %v\n", filename, err)
			continue
		}
		
		// Write auction content
		fmt.Fprintf(file, "=== Auction %d Results ===\n", result.AuctionID)
		fmt.Fprintf(file, "Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
		
		if result.WinnerID == -1 {
			fmt.Fprintf(file, "Winner: No valid bids received\n")
			fmt.Fprintf(file, "Winning Amount: N/A\n")
		} else {
			fmt.Fprintf(file, "Winner: Bidder %d\n", result.WinnerID)
			fmt.Fprintf(file, "Winning Amount: %.2f\n", result.Amount)
		}
		
		fmt.Fprintf(file, "Total Bids Received: %d\n\n", len(result.Bids))
		
		// Write all bids
		if len(result.Bids) > 0 {
			fmt.Fprintf(file, "=== All Bids ===\n")
			fmt.Fprintf(file, "%-10s | %-10s\n", "Bidder ID", "Amount")
			fmt.Fprintf(file, "----------------------\n")
			
			// Sort bids by amount
			bids := make([]auction.Bid, len(result.Bids))
			copy(bids, result.Bids)
			sort.Slice(bids, func(i, j int) bool { return bids[i].Amount > bids[j].Amount })
			
			for _, bid := range bids {
				fmt.Fprintf(file, "%-10d | %-10.2f\n", bid.BidderID, bid.Amount)
			}
		}
		
		file.Close()
		successCount++
	}
	
	fmt.Printf("Created %d individual auction files\n", successCount)
}