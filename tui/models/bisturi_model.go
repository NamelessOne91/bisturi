package tui

import (
	"fmt"
	"net"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step uint8

const (
	setupStep step = iota
	ifaceStep
	protoStep
	doneStep
)

type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }

type bisturiModel struct {
	step              step
	startMenu         startMenuModel
	spinner           spinner.Model
	selectedInterface *net.Interface
	selectedProtocol  string
	selectedEthType   uint16
	err               error
}

func NewBisturiModel() *bisturiModel {
	s := spinner.New(spinner.WithSpinner(spinner.Meter))
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99"))

	return &bisturiModel{
		step:      setupStep,
		spinner:   s,
		startMenu: startMenuModel{},
	}
}

func (m bisturiModel) Init() tea.Cmd {
	return tea.Sequence(m.spinner.Tick, fetchInterfaces())
}

func (m *bisturiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.err
		return m, tea.Quit

	case networkInterfacesMsg:
		m.startMenu = newStartMenuModel(msg)
		m.step = ifaceStep
		return m, nil

	case selectedInterfaceMsg:
		iface, err := net.InterfaceByName(msg.name)
		if err != nil {
			return m, func() tea.Msg {
				return errMsg{err: err}
			}
		}
		m.selectedInterface = iface
		m.step = protoStep
		m.startMenu.step = protoStep
		return m, nil

	case selectedProtocolMsg:
		m.selectedProtocol = msg.protocol
		m.selectedEthType = msg.ethTytpe
		m.step = doneStep
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.startMenu, cmd = m.startMenu.Update(msg)
	return m, cmd
}

func (m bisturiModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %s\n", m.err)
	}

	switch m.step {
	case ifaceStep, protoStep:
		return m.startMenu.View()
	case doneStep:
		return fmt.Sprintf("Receiving packet on %s - %s\n", m.selectedInterface.Name, m.selectedProtocol)
	default:
		return fmt.Sprintf("\n\n %s Retrieving network interfaces...\n\n", m.spinner.View())
	}

}
