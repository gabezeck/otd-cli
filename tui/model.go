package tui

import (
	"fmt"
	"strings"

	"otd-cli/scraper"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	spinner  spinner.Model
	viewport viewport.Model
	data     *scraper.OTDData
	err      error
	ready    bool
	width    int
	height   int
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchData)
}

func fetchData() tea.Msg {
	data, err := scraper.Scrape()
	if err != nil {
		return errMsg(err)
	}
	return data
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Header height approx 2-3 lines, Footer 1 line
		// We'll give the viewport the rest
		headerHeight := 2
		footerHeight := 1
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.viewport.SetContent(m.renderContent()) // Re-render with new width
		}

	case errMsg:
		m.err = msg
		return m, nil

	case *scraper.OTDData:
		m.data = msg
		m.ready = true
		// Set content now that we have data
		if m.width > 0 {
			// Initialize viewport if it hasn't been yet (rare race case, but possible)
			if m.viewport.Height == 0 {
				headerHeight := 2
				footerHeight := 1
				verticalMarginHeight := headerHeight + footerHeight
				m.viewport = viewport.New(m.width, m.height-verticalMarginHeight)
				m.viewport.YPosition = headerHeight
			}
			m.viewport.SetContent(m.renderContent())
		}

	case spinner.TickMsg:
		if !m.ready {
			var newSpinnerCmd tea.Cmd
			m.spinner, newSpinnerCmd = m.spinner.Update(msg)
			cmds = append(cmds, newSpinnerCmd)
		}
	}

	// Only update viewport if ready
	if m.ready {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) renderContent() string {
	if m.data == nil {
		return ""
	}

	var sb strings.Builder

	// Events Section
	sb.WriteString(sectionHeaderStyle.Render("Historical Events"))
	sb.WriteString("\n")

	for _, e := range m.data.Events {
		// Layout: [Year] │ [Description]

		year := yearStyle.Render(e.Year)
		dot := timelineDotStyle.String()

		// Calculate available width for description
		// width - year(6) - margin(1) - dot(1) - margin(1) = width - 9
		descWidth := m.width - 12
		if descWidth < 20 {
			descWidth = 20
		}

		desc := descStyle.Width(descWidth).Render(e.Text)

		row := lipgloss.JoinHorizontal(lipgloss.Top, year, dot, desc)
		sb.WriteString(row)
		sb.WriteString("\n\n")
	}

	// Birthdays Section
	sb.WriteString("\n")
	sb.WriteString(sectionHeaderStyle.Render("Famous Birthdays"))
	sb.WriteString("\n")

	for _, b := range m.data.Birthdays {
		name := nameStyle.Render(b.Name)
		year := bYearStyle.Render(b.Year)
		line := fmt.Sprintf("• %s%s", name, year)
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\nPress q to quit.", m.err)
	}

	if !m.ready {
		return fmt.Sprintf("\n %s Loading history...\n", m.spinner.View())
	}

	// Header
	// "ON THIS DAY" | "December 21"
	title := headerTitleStyle.Render("ON THIS DAY")
	date := headerDateStyle.Render(m.data.Date)
	header := lipgloss.JoinHorizontal(lipgloss.Center, title, date)
	header = lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(header)

	// Viewport
	content := m.viewport.View()

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(subtle).
		Align(lipgloss.Center).
		Width(m.width).
		Render("Scroll: j/k • Quit: q")

	return fmt.Sprintf("%s\n\n%s\n%s", header, content, footer)
}
