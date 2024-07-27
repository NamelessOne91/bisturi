package tui

import (
	"fmt"
	"time"

	"github.com/NamelessOne91/bisturi/sockets"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyID          = "id"
	columnKeyDate        = "date"
	columnKeySource      = "source"
	columnKeyDestination = "destination"
	columnKeyInfo        = "info"
)

type packetsTablemodel struct {
	table      table.Model
	maxRows    int
	cachedRows []table.Row
	counter    uint64
}

func buildTable(rows []table.Row, terminalWidth int) table.Model {
	return table.New([]table.Column{
		table.NewColumn(columnKeyID, "#", (2*terminalWidth)/100),
		table.NewColumn(columnKeyDate, "Date", (8*terminalWidth)/100),
		table.NewColumn(columnKeySource, "Source", (20*terminalWidth)/100),
		table.NewColumn(columnKeyDestination, "Destination", (20*terminalWidth)/100),
		table.NewColumn(columnKeyInfo, "Info", (46*terminalWidth)/100),
	}).
		WithRows(rows).
		WithBaseStyle(lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("#00cc99")).
			Foreground(lipgloss.Color("#00cc99")).
			Align(lipgloss.Center),
		)
}

func newPacketsTable(max int, terminalWidth int) packetsTablemodel {
	rows := make([]table.Row, 0, max)

	return packetsTablemodel{
		maxRows:    max,
		cachedRows: rows,
		table:      buildTable(rows, terminalWidth),
	}
}

func (m *packetsTablemodel) resizeTable(terminalWidth int) {
	m.table = buildTable(m.cachedRows, terminalWidth)
}

func (m packetsTablemodel) Init() tea.Cmd {
	return nil
}

func (m packetsTablemodel) Update(msg tea.Msg) (packetsTablemodel, tea.Cmd) {
	switch msg := msg.(type) {
	case readPacketsMsg:
		m.addRows(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m packetsTablemodel) View() string {
	return fmt.Sprintf("Displaying up to the last %d rows\n\n%s", m.maxRows, m.table.View())
}

func (m *packetsTablemodel) addRows(packets []sockets.NetworkPacket) {
	for _, np := range packets {
		if len(m.cachedRows) >= m.maxRows {
			m.cachedRows = m.cachedRows[1:]
		}
		m.counter += 1

		newRow := table.NewRow(table.RowData{
			columnKeyID:          m.counter,
			columnKeyDate:        time.Now().Local().Format(time.Stamp),
			columnKeySource:      np.Source(),
			columnKeyDestination: np.Destination(),
			columnKeyInfo:        np.Info(),
		})
		m.cachedRows = append(m.cachedRows, newRow)
	}
	m.table = m.table.WithRows(m.cachedRows)
}
