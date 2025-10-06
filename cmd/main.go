package main

import (
	"fmt"
	"log/slog"
	"os"
	"wobitoxmas/internal/draw"

	"github.com/olekukonko/tablewriter"
)

func main() {
	d, err := draw.NewDraw()
	if err != nil {
		slog.Error("failed to initialize draw", slog.String("error", err.Error()))
		os.Exit(1)
	}

	fmt.Println("Running Christmas Gift Draw...")
	results, err := d.StartDraw()
	if err != nil {
		slog.Error("failed to complete draw", slog.String("error", err.Error()))
		fmt.Println("\nNo valid assignment possible with current constraints.")
		fmt.Println("Try adjusting the LastYearID values.")
		os.Exit(1)
	}

	// Convert results to table format
	tableData := make([][]string, 0, len(results))
	for _, result := range results {
		tableData = append(tableData, []string{result.Giver, result.Receiver})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Name", "Picked"})
	table.Bulk(tableData)
	table.Render()
}
