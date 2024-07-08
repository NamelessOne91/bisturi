package tui

import (
	"net"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ifaceItem represents a network interface in a list
type ifaceItem struct {
	name  string
	flags string
}

func (i ifaceItem) Title() string { return i.name }

func (i ifaceItem) Description() string { return i.flags }

func (i ifaceItem) FilterValue() string { return i.name }

type selectedIfaceItemMsg ifaceItem

type networkInterfacesMsg []net.Interface

type interfacesListModel struct {
	l list.Model
}

func (m interfacesListModel) Init() tea.Cmd {
	return nil
}

func (m interfacesListModel) Update(msg tea.Msg) (interfacesListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.l.SelectedItem().(ifaceItem)
			if ok {
				return m, func() tea.Msg {
					return selectedIfaceItemMsg(i)
				}
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.l, cmd = m.l.Update(msg)
	return m, cmd
}

func (m interfacesListModel) View() string {
	return m.l.View()
}

func newInterfacesListModel(width, height int, interfaces []net.Interface) interfacesListModel {
	items := make([]list.Item, len(interfaces))
	for i, iface := range interfaces {
		items[i] = ifaceItem{
			name:  iface.Name,
			flags: iface.Flags.String(),
		}
	}

	titlesStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99")).Blink(true).Bold(true)
	ifaceDelegate := list.NewDefaultDelegate()
	ifaceDelegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00cc99"))

	ifaceList := list.New(items, ifaceDelegate, width, height)
	ifaceList.Title = "Select a Network Interface"
	ifaceList.Styles.Title = titlesStyle
	ifaceList.SetShowStatusBar(true)
	ifaceList.SetFilteringEnabled(false)
	ifaceList.SetShowHelp(true)

	return interfacesListModel{l: ifaceList}
}

func fetchInterfaces() tea.Cmd {
	return func() tea.Msg {
		ifaces, err := net.Interfaces()
		if err != nil {
			return errMsg{err: err}
		}
		return networkInterfacesMsg(ifaces)
	}
}
