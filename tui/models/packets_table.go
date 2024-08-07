package tui

import (
	"time"

	"github.com/NamelessOne91/bisturi/sockets"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	columnKeyID          = "id"
	columnKeyTime        = "time"
	columnKeySource      = "source"
	columnKeyDestination = "destination"
	columnKeyInfo        = "info"
)

type packetsTableModel struct {
	table      table.Model
	height     int
	width      int
	maxRows    int
	cachedRows []table.Row
	counter    uint64
}

func (m *packetsTableModel) buildTable() {
	m.table = table.New([]table.Column{
		table.NewColumn(columnKeyID, "#", (3*m.width)/100),
		table.NewColumn(columnKeyTime, "Time", (7*m.width)/100),
		table.NewColumn(columnKeySource, "Source", (20*m.width)/100),
		table.NewColumn(columnKeyDestination, "Destination", (20*m.width)/100),
	}).
		WithRows(m.cachedRows).
		Focused(true).
		WithBaseStyle(lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("#00cc99")).
			Foreground(lipgloss.Color("#00cc99")).
			Align(lipgloss.Center),
		)
}

func newPacketsTable(max int, height, width int) packetsTableModel {
	rows := make([]table.Row, 0, max)

	ptm := packetsTableModel{
		height:     height,
		width:      width,
		maxRows:    max,
		cachedRows: rows,
	}
	ptm.buildTable()

	return ptm
}

func (m *packetsTableModel) resize(height, width int) {
	m.height = height
	m.width = width
	m.buildTable()
}

func (m packetsTableModel) Init() tea.Cmd {
	return nil
}

func (m packetsTableModel) Update(msg tea.Msg) (packetsTableModel, tea.Cmd) {
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

func (m packetsTableModel) View() string {
	var detailTxt string
	if len(m.table.GetVisibleRows()) > 0 {
		detailTxt = m.table.HighlightedRow().Data[columnKeyInfo].(string)
	}

	detailsBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#00cc99")).
		Foreground(lipgloss.Color("#00cc99")).
		Padding(1, 2).
		Width((40 * m.width) / 100).
		Height((90 * m.height) / 100).
		Render(detailTxt)

	mainView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.table.View(),
		detailsBox,
	)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		mainView,
	) + "\n"

	return view
}

func (m *packetsTableModel) addRows(packets []sockets.NetworkPacket) {
	lp := len(packets)
	lc := len(m.cachedRows)

	if (lp + lc) > m.maxRows {
		newCache := make([]table.Row, 0, m.maxRows)
		if lp > m.maxRows {
			packets = packets[lp-m.maxRows:]
		} else {
			oldToKeep := m.maxRows - lp
			newCache = append(newCache, m.cachedRows[lc-oldToKeep:]...)
		}
		m.cachedRows = newCache
	}

	for _, np := range packets {
		m.counter += 1

		newRow := table.NewRow(table.RowData{
			columnKeyID:          m.counter,
			columnKeyTime:        time.Now().Local().Format(time.TimeOnly),
			columnKeySource:      np.Source(),
			columnKeyDestination: np.Destination(),
			columnKeyInfo:        np.Info(),
		})
		m.cachedRows = append(m.cachedRows, newRow)
	}
	m.table = m.table.WithRows(m.cachedRows)
}
