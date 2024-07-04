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

type interfacesListModel struct {
	l list.Model
}

func newInterfacesListModel(width, height int) interfacesListModel {
	items := []list.Item{}

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
