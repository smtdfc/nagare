package router

import tea "charm.land/bubbletea/v2"

type Page interface {
	tea.Model
	GetName() string
	Refresh()
}

type TUIRouter struct {
	Pages       map[string]Page
	CurrentPage string
}

func (r *TUIRouter) HasPage(name string) bool {
	_, ok := r.Pages[name]
	return ok
}

func NewTUIRouter() *TUIRouter {
	return &TUIRouter{
		Pages:       map[string]Page{},
		CurrentPage: "main",
	}
}
