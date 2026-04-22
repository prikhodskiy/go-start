package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.front == nil {
		// List is empty
		l.front = newItem
		l.back = newItem
	} else {
		// Insert at front
		newItem.Next = l.front
		l.front.Prev = newItem
		l.front = newItem
	}

	l.len++
	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.back == nil {
		// List is empty
		l.front = newItem
		l.back = newItem
	} else {
		// Insert at back
		newItem.Prev = l.back
		l.back.Next = newItem
		l.back = newItem
	}

	l.len++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		// i is the front element
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		// i is the back element
		l.back = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.front {
		return
	}

	// Remove i from its current position
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		// i is the back element
		l.back = i.Prev
	}

	// Insert i at front
	i.Next = l.front
	i.Prev = nil
	l.front.Prev = i
	l.front = i
}
