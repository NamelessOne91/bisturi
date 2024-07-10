package tui

import (
	"fmt"
	"net"
	"strings"

	"github.com/NamelessOne91/bisturi/sockets"
	"github.com/NamelessOne91/bisturi/tui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step uint8

const (
	retrieveIfaces step = iota
	selectIface
	selectProtocol
	initPackets
	receivePackets
)

type errMsg error

type packetMsg sockets.NetworkPacket

type bisturiModel struct {
	step              step
	spinner           spinner.Model
	startMenu         startMenuModel
	packetsTable      packetsTablemodel
	selectedInterface net.Interface
	selectedProtocol  string
	selectedEthType   uint16
	rawSocket         *sockets.RawSocket
	packetsChan       chan sockets.NetworkPacket
	errChan           chan error
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
	case initPackets:
		m.step = receivePackets
		return m, m.waitForPacket()
	case receivePackets:
		return m.updateReceivingPacket(msg)
	default:
		return m, nil
	}
}

func (m bisturiModel) View() string {
	sb := strings.Builder{}

	if m.err != nil {
		sb.WriteString(fmt.Sprintf("Error: %s\n", m.err))
		if m.rawSocket != nil {
			if err := m.rawSocket.Close(); err != nil {
				sb.WriteString(err.Error())
			}
		}
	}

	switch m.step {
	case retrieveIfaces:
		sb.WriteString(fmt.Sprintf("\n\nWelcome!\n\n %s Retrieving network interfaces...\n\n", m.spinner.View()))
	case selectIface, selectProtocol:
		sb.WriteString(m.startMenu.View())
	case initPackets, receivePackets:
		sb.WriteString(fmt.Sprintf("\n\nReceiving %s packets on %s ...\n\n", m.selectedProtocol, m.selectedInterface.Name))
		sb.WriteString(m.packetsTable.View())
	default:
		sb.WriteString("The program is in an unknown state\nQuit with 'q'")
	}
	return styles.Default.Render(sb.String())
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

		err = rs.Bind(m.selectedInterface)
		if err != nil {
			m.err = err
			return m, tea.Quit
		}
		m.selectedProtocol = msg.name
		m.selectedEthType = msg.ethType
		m.rawSocket = rs
		m.step = initPackets
		m.packetsTable = newPacketsTable()

		m.packetsChan = make(chan sockets.NetworkPacket)
		m.errChan = make(chan error)

		return m, func() tea.Msg {
			go m.rawSocket.ReadToChan(m.packetsChan, m.errChan)
			return nil
		}
	}
	return m, cmd
}

func (m *bisturiModel) updateReceivingPacket(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.packetsTable, cmd = m.packetsTable.Update(msg)

	switch msg.(type) {
	case sockets.NetworkPacket:
		return m, m.waitForPacket()
	}
	return m, cmd
}

func (m bisturiModel) waitForPacket() tea.Cmd {
	return func() tea.Msg {
		select {
		case packet := <-m.packetsChan:
			return packetMsg(packet)
		case err := <-m.errChan:
			return errMsg(err)
		}
	}
}
