package image_cache //nolint

type List struct {
	front *ListItem
	back  *ListItem
	len   int
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

func NewList() *List {
	return &List{}
}

func (l List) Len() int {
	return l.len
}

func (l List) Front() *ListItem {
	return l.front
}

func (l List) Back() *ListItem {
	return l.back
}

func (l *List) PushFront(v interface{}) *ListItem {
	front := l.Front()

	newItem := &ListItem{
		Value: v,
		Next:  front,
		Prev:  nil,
	}

	if front != nil {
		front.Prev = newItem
	} else {
		l.back = newItem
	}
	l.len++
	l.front = newItem

	return newItem
}

func (l *List) PushBack(v interface{}) *ListItem {
	back := l.Back()

	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  back,
	}

	if back != nil {
		back.Next = newItem
	} else {
		l.front = newItem
	}
	l.len++
	l.back = newItem

	return newItem
}

func (l *List) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	l.len--
}

func (l *List) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
