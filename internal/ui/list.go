package ui

import "github.com/charmbracelet/bubbles/list"

type item struct {
	title       string
	description string
}

func (i item) Title() string {
	return i.title
}

func (i item) Description() string {
	return i.description
}

func (i item) FilterValue() string {
	return i.title
}

type menu struct {
	list list.Model
}

func (m menu) New() *menu {
	m.list = list.New(
		[]list.Item{item{"Games", "Game List"}, item{"Standings", "Standing"}},
		list.NewDefaultDelegate(),
		0, 0)

	return &m
}
