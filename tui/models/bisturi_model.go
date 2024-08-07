package tui

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/NamelessOne91/bisturi/sockets"
	"github.com/NamelessOne91/bisturi/tui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type step uint8

const (
	retrieveIfaces step = iota
	selectIface
	selectProtocol
	selectRows
	receivePackets
)

type errMsg error

type readPacketsMsg []sockets.NetworkPacket

type bisturiModel struct {
	terminalHeight    int
	terminalWidth     int
	step              step
	spinner           spinner.Model
	startMenu         startMenuModel
	rowsInput         textinput.Model
	packetsTable      packetsTableModel
	selectedInterface net.Interface
	selectedProtocol  string
	selectedEthType   uint16
	rawSocket         *sockets.RawSocket
	packetsChan       chan sockets.NetworkPacket
	msgChan           chan tea.Msg
	errChan           chan error
	err               error
}

func newRowsInput(terminalWidth int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Enter the max number of rows to display"
	ti.Focus()
	ti.CharLimit = 4
	ti.Width = terminalWidth / 2

	return ti
}

func NewBisturiModel() *bisturiModel {
	s := spinner.New(spinner.WithSpinner(spinner.Meter))
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99"))

	return &bisturiModel{
		step:        retrieveIfaces,
		spinner:     s,
		packetsChan: make(chan sockets.NetworkPacket),
		msgChan:     make(chan tea.Msg),
		errChan:     make(chan error),
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
	case selectRows:
		return m.updateRowsInput(msg)
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
		sb.WriteString(fmt.Sprintf("\nWelcome!\nRetrieving network interfaces \n\n%s", m.spinner.View()))
	case selectIface, selectProtocol:
		sb.WriteString(m.startMenu.View())
	case selectRows:
		sb.WriteString(m.rowsInput.View())
	case receivePackets:
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

	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width

		return m, nil

	case networkInterfacesMsg:
		m.startMenu = newStartMenuModel(msg, m.terminalHeight, m.terminalWidth)
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
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Height
		if m.step == selectIface {
			m.startMenu.ifaceList.resize(m.terminalHeight, m.terminalWidth)
		} else if m.step == selectProtocol {
			m.startMenu.protoList.resize(m.terminalHeight, m.terminalWidth)
		}

		return m, nil

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
			m.err = err
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
		m.step = selectRows

		m.rowsInput = newRowsInput(m.terminalWidth)
		return m, nil
	}
	return m, cmd
}

func (m *bisturiModel) updateRowsInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.rowsInput, cmd = m.rowsInput.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit

		case "enter":
			maxRows, err := strconv.Atoi(m.rowsInput.Value())
			if err == nil && maxRows > 0 {
				m.packetsTable = newPacketsTable(maxRows, m.terminalHeight, m.terminalWidth)
				m.step = receivePackets

				go m.rawSocket.ReadToChan(m.packetsChan, m.errChan)
				go m.readPackets()

				return m, m.pollPacketsMessages()
			}
		}
	}

	return m, cmd
}

func (m *bisturiModel) updateReceivingPacket(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.packetsTable, cmd = m.packetsTable.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.terminalHeight = msg.Height
		m.terminalWidth = msg.Width
		m.packetsTable.resize(m.terminalHeight, m.terminalWidth)

		return m, nil
	case readPacketsMsg:
		return m, m.pollPacketsMessages()
	}
	return m, cmd
}

func (m bisturiModel) readPackets() {
	readPackets := []sockets.NetworkPacket{}
	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()

	for {
		select {
		case packet := <-m.packetsChan:
			readPackets = append(readPackets, packet)
		case <-timer.C:
			if len(readPackets) > 0 {
				m.msgChan <- readPacketsMsg(readPackets)
				readPackets = []sockets.NetworkPacket{}
			}
		case err := <-m.errChan:
			m.msgChan <- errMsg(err)
		}
	}
}

func (m bisturiModel) pollPacketsMessages() tea.Cmd {
	return func() tea.Msg {
		return <-m.msgChan
	}
}
