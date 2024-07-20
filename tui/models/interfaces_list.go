package tui

import (
	"net"
	"time"

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

func newInterfacesListModel(interfaces []net.Interface, terminalHeight, terminalWidth int) interfacesListModel {
	listHeight := (75 * terminalHeight) / 100
	listWidth := (75 * terminalWidth) / 100

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

	ifaceList := list.New(items, ifaceDelegate, listWidth, listHeight)
	ifaceList.Title = "Select a Network Interface"
	ifaceList.Styles.Title = titlesStyle
	ifaceList.SetShowStatusBar(true)
	ifaceList.SetFilteringEnabled(false)
	ifaceList.SetShowHelp(true)

	return interfacesListModel{l: ifaceList}
}

func (m *interfacesListModel) resize(terminalHeight, terminalWidth int) {
	listHeight := (75 * terminalHeight) / 100
	listWidth := (75 * terminalWidth) / 100

	m.l.SetHeight(listHeight)
	m.l.SetWidth(listWidth)
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

func fetchInterfaces() tea.Cmd {
	return func() tea.Msg {
		ifaces, err := net.Interfaces()
		if err != nil {
			return errMsg(err)
		}
		time.Sleep(2 * time.Second)
		return networkInterfacesMsg(ifaces)
	}
}
