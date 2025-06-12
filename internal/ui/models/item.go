package models

// Item represents a list item
type Item struct {
	title       string
	description string
	Action      string
}

func NewItem(title, description, action string) Item {
	return Item{title: title, description: description, Action: action}
}

func (i Item) FilterValue() string { return i.title }
func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.description }