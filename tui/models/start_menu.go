package tui

import (
	"fmt"
	"net"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type startMenuStep uint

const (
	ifaceStep startMenuStep = iota
	protoStep
	recvStep
)

type startMenuModel struct {
	step           startMenuStep
	ifaceListModel interfacesListModel
	protoListModel protocolsListModel
	chosenIface    *net.Interface
	chosenProto    string
	chosenEthType  uint16
	err            error
}

func NewStartMenuModel() *startMenuModel {
	const listHeight = 50
	const listWidth = 50

	ilm := newInterfacesListModel(listWidth, listHeight)
	plm := newProtocolsListModel(listWidth, listHeight)

	return &startMenuModel{
		ifaceListModel: ilm,
		protoListModel: plm,
	}
}

func (m startMenuModel) Init() tea.Cmd {
	return fetchInterfaces()
}

func (m startMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case ifaceStep:
		return m.updateInterfacesListModel(msg)
	case protoStep:
		return m.updateProtocolsListModel(msg)
	case recvStep:
		return m, tea.Quit
	}
	return m, nil
}

func (m startMenuModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %s\n", m.err)
	}

	var s string
	switch m.step {
	case ifaceStep:
		s = m.ifaceListModel.l.View()
	case protoStep:
		s = m.protoListModel.l.View()
	case recvStep:
		s = fmt.Sprintf("You chose %s - %s\n", m.chosenIface.Name, m.chosenProto)
	default:
		return "Unkown step"
	}
	return lipgloss.NewStyle().Padding(1).Render(s)
}

func (m *startMenuModel) updateInterfacesListModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg.err
		return m, tea.Quit

	case networkInterfacesMsg:
		items := make([]list.Item, len(msg))
		for i, iface := range msg {
			items[i] = ifaceItem{
				name:  iface.Name,
				flags: iface.Flags.String(),
			}
		}
		cmd := m.ifaceListModel.l.SetItems(items)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			item, ok := m.ifaceListModel.l.SelectedItem().(ifaceItem)
			if ok {
				// retrieve the network interface
				networkInterface, err := net.InterfaceByName(item.name)
				if err != nil {
					m.err = err
					return m, tea.Quit
				}
				m.chosenIface = networkInterface
				m.step = protoStep
			}
		}
	}

	var cmd tea.Cmd
	m.ifaceListModel.l, cmd = m.ifaceListModel.l.Update(msg)
	return m, cmd
}

func (m *startMenuModel) updateProtocolsListModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			p, ok := m.protoListModel.l.SelectedItem().(protoItem)
			if ok {
				m.chosenProto = p.name
				m.chosenEthType = p.ethType
				m.step = recvStep
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.protoListModel.l, cmd = m.protoListModel.l.Update(msg)
	return m, cmd
}
