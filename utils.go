package main

import "fmt"

func prettifySize(b int64) string {
	// Unit size
	const unit = 1024

	// If the input is less than 1KB, show bytes suffix
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	// Convert unit size to int64, initialize counter variable
	div, exp := int64(unit), 0

	// Keep dividing the input value by the unit size until it equals less than the unit size
	for n := b / unit; n >= unit; n /= unit {
		// Keep track of the total units we've divided by
		div *= unit

		// Increment the counter
		exp++
	}

	// Divide the input value by the final count of units
	// Use the counter for number of times divided to get the unit suffix
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
