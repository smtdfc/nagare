package components

import "fmt"

type ListItem[T any] struct {
	title string
	desc  string
	value T
}

func (i ListItem[T]) Title() string       { return i.title }
func (i ListItem[T]) Description() string { return i.desc }
func (i ListItem[T]) Value() T            { return i.value }
func (i ListItem[T]) FilterValue() string { return fmt.Sprint(i.value) }

func NewListItem[T any](title, desc string, value T) ListItem[T] {
	return ListItem[T]{
		title: title,
		desc:  desc,
		value: value,
	}
}
