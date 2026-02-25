package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"otd-cli/scraper"
	"otd-cli/tui"
)

func main() {
	headless := flag.Bool("headless", false, "Print output and exit (no TUI)")
	flag.Parse()

	if *headless {
		data, err := scraper.Scrape()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Alas, there's been an error: %v\n", err)
			os.Exit(1)
		}

		output := tui.RenderHeadless(data, terminalWidth())
		fmt.Println(output)
		return
	}

	p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func terminalWidth() int {
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if width, err := strconv.Atoi(cols); err == nil && width > 0 {
			return width
		}
	}

	return 80
}
