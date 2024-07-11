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

func newPacketsTable(max int) packetsTablemodel {
	rows := make([]table.Row, 0, max)

	return packetsTablemodel{
		maxRows:    max,
		cachedRows: rows,
		table: table.New([]table.Column{
			table.NewColumn(columnKeyID, "#", 5),
			table.NewColumn(columnKeyDate, "Date", 20),
			table.NewColumn(columnKeySource, "Source", 30),
			table.NewColumn(columnKeyDestination, "Destination", 30),
			table.NewColumn(columnKeyInfo, "Info", 100),
		}).
			WithRows(rows).
			WithBaseStyle(lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#00cc99")).
				Foreground(lipgloss.Color("#00cc99")).
				Align(lipgloss.Center),
			),
	}
}

func (m packetsTablemodel) Init() tea.Cmd {
	return nil
}

func (m packetsTablemodel) Update(msg tea.Msg) (packetsTablemodel, tea.Cmd) {
	switch msg := msg.(type) {
	case packetMsg:
		m.addRow(msg)
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

func (m *packetsTablemodel) addRow(data sockets.NetworkPacket) {
	if len(m.cachedRows) >= m.maxRows {
		m.cachedRows = m.cachedRows[1:]
	}
	m.counter += 1

	newRow := table.NewRow(table.RowData{
		columnKeyID:          m.counter,
		columnKeyDate:        time.Now().Local().Format(time.Stamp),
		columnKeySource:      data.Source(),
		columnKeyDestination: data.Destination(),
		columnKeyInfo:        data.Info(),
	})
	m.cachedRows = append(m.cachedRows, newRow)
	m.table = m.table.WithRows(m.cachedRows)
}
