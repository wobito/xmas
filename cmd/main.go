package main

import (
	"fmt"
	"os"
	"wobitoxmas/internal/draw"

	"github.com/olekukonko/tablewriter"
)

func main() {
	d := draw.NewDraw()
	d.SetUsers()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Picked"})
	var i int
	for i < 20 {
		i++
		fmt.Printf("Running Draw #%d ....\n", i)
		for table.NumLines() != len(d.Users) {
			table.ClearRows()
			d.StartDraw(table)
			if table.NumLines() != len(d.Users) {
				d.SetUsers()
			}
		}
	}

	table.Render() // Send output
}
