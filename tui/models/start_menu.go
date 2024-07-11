package tui

import (
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

type startMenuModel struct {
	step      step
	ifaceList interfacesListModel
	protoList protocolsListModel
}

type selectedInterfaceMsg struct {
	name string
}

type selectedProtocolMsg struct {
	protocol string
	ethTytpe uint16
}

func newStartMenuModel(interfaces []net.Interface) startMenuModel {
	const listHeight = 50
	const listWidth = 50

	il := newInterfacesListModel(listWidth, listHeight, interfaces)
	plm := newProtocolsListModel(listWidth, listHeight)

	return startMenuModel{
		step:      selectIface,
		ifaceList: il,
		protoList: plm,
	}
}

func (m startMenuModel) Init() tea.Cmd {
	return nil
}

func (m startMenuModel) Update(msg tea.Msg) (startMenuModel, tea.Cmd) {
	switch m.step {
	case selectIface:
		model, cmd := m.ifaceList.Update(msg)
		if i, ok := msg.(selectedIfaceItemMsg); ok {
			m.step = selectProtocol
			return m, func() tea.Msg {
				return selectedInterfaceMsg{
					name: i.name,
				}
			}
		}
		m.ifaceList = model
		return m, cmd

	case selectProtocol:
		model, cmd := m.protoList.Update(msg)
		if p, ok := msg.(selectedProtocolItemMsg); ok {
			return m, func() tea.Msg {
				return selectedProtocolMsg{
					protocol: p.name,
					ethTytpe: p.ethType,
				}
			}
		}
		m.protoList = model
		return m, cmd
	}
	return m, nil
}

func (m startMenuModel) View() string {
	var s string
	switch m.step {
	case selectIface:
		s = m.ifaceList.l.View()
	case selectProtocol:
		s = m.protoList.l.View()
	default:
		s = "Unkown step"
	}
	return s
}
