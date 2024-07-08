package tui

import (
	"fmt"
	"net"
	"strings"

	"github.com/NamelessOne91/bisturi/sockets"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step uint8

const (
	retrieveIfaces step = iota
	selectIface
	selectProtocol
	receivePackets
)

type errMsg error

type bisturiModel struct {
	step              step
	spinner           spinner.Model
	startMenu         startMenuModel
	packetsTable      packetsTablemodel
	selectedInterface net.Interface
	selectedProtocol  string
	selectedEthType   uint16
	rawSocket         *sockets.RawSocket
	err               error
}

func NewBisturiModel() *bisturiModel {
	s := spinner.New(spinner.WithSpinner(spinner.Meter))
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99"))

	return &bisturiModel{
		step:    retrieveIfaces,
		spinner: s,
	}
}

func (m bisturiModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchInterfaces())
}

func (m bisturiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case retrieveIfaces:
		return m.updateLoading(msg)
	case selectIface, selectProtocol:
		return m.updateStartMenuSelection(msg)
	case receivePackets:
		return m.updateReceivingPacket(msg)
	default:
		return m, nil
	}
}

func (m *bisturiModel) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, tea.Quit

	case networkInterfacesMsg:
		m.startMenu = newStartMenuModel(msg)
		m.step = selectIface

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m *bisturiModel) updateStartMenuSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.startMenu, cmd = m.startMenu.Update(msg)

	switch msg := msg.(type) {
	case selectedIfaceItemMsg:
		iface, err := net.InterfaceByName(msg.name)
		if err != nil {
			m.err = err
			return m, tea.Quit
		}
		m.selectedInterface = *iface
		m.step = selectProtocol

		return m, nil

	case selectedProtocolItemMsg:
		// SYS_SOCKET syscall
		rs, err := sockets.NewRawSocket(msg.name, msg.ethType)
		if err != nil {
			return m, tea.Quit
		}
		// bind the socket to the network interface
		err = rs.Bind(m.selectedInterface)
		if err != nil {
			m.err = err
			return m, tea.Quit
		}
		m.selectedProtocol = msg.name
		m.selectedEthType = msg.ethType
		m.rawSocket = rs
		m.step = receivePackets
		m.packetsTable = newPacketsTable()

		return m, nil
	}
	return m, cmd
}

func (m *bisturiModel) updateReceivingPacket(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.packetsTable, cmd = m.packetsTable.Update(msg)

	return m, cmd
}

func (m bisturiModel) View() string {
	if m.err != nil {
		if m.rawSocket != nil {
			m.rawSocket.Close()
		}
		return fmt.Sprintf("Error: %s\n", m.err)
	}

	sb := strings.Builder{}
	switch m.step {
	case retrieveIfaces:
		sb.WriteString(fmt.Sprintf("\n\n %s Retrieving network interfaces...\n\n", m.spinner.View()))
	case selectIface, selectProtocol:
		sb.WriteString(m.startMenu.View())
	case receivePackets:
		sb.WriteString(fmt.Sprintf("Receiving %s packets on %s ...\n", m.selectedProtocol, m.selectedInterface.Name))
		sb.WriteString(m.packetsTable.View())
	default:
		sb.WriteString("The program is in an unkowqn state\n")
	}
	return sb.String()
}
