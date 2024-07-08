package tui

import (
	"strings"

	"github.com/NamelessOne91/bisturi/protocols"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyProtocol  = "protocol"
	columnKeyInterface = "interface"
	columnKeyInfo      = "info"
)

type packetsTablemodel struct {
	table       table.Model
	cachedRows  []table.Row
	packetsChan <-chan protocols.IPPacket
}

func newPacketsTable() packetsTablemodel {
	rows := make([]table.Row, 0, 20)

	return packetsTablemodel{
		cachedRows: rows,
		table: table.New([]table.Column{
			table.NewColumn(columnKeyInterface, "Interface", 20),
			table.NewColumn(columnKeyProtocol, "Protocol", 20),
			table.NewColumn(columnKeyProtocol, "Info", 50),
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
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m packetsTablemodel) View() string {
	sb := strings.Builder{}

	sb.WriteString(m.table.View())

	return sb.String()
}
