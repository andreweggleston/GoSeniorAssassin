package database

import "fmt"

// A linkedlist implementation which does not allow empty lists
// empty lists are basically just nil values
type LinkedList struct {
	value interface{}
	next  *LinkedList
}

func NewLinkedList(values []interface{}) *LinkedList {
	if len(values) == 0 {
		return nil
	}
	first := &LinkedList{
		value: values[0],
		next:  nil,
	}
	for _, value := range values[1:] {
		first.Push(value)
	}
	return first
}

func (l *LinkedList) Len() int {
	if l == nil {
		return 0
	}
	if l.next == nil {
		return 1
	}
	return 1 + l.next.Len()
}

func (l *LinkedList) Push(val interface{}) (*LinkedList, error) {
	return l.Insert(l.Len()-1, val)
}

func (l *LinkedList) Insert(index int, val interface{}) (*LinkedList, error) {
	if index == 0 {
		insertion := &LinkedList{
			value: val,
			next:  l.next,
		}
		l.next = insertion
		return l, nil
	} else {
		if l.next == nil {
			return nil, fmt.Errorf("Out of bounds")
		} else {
			l.next.Insert(index-1, val)
			return l, nil
		}
	}
}

func (l *LinkedList) Find(val interface{}) (int, error) {
	if l.value == val {
		return 0, nil
	}
	if l.next == nil {
		return -1, fmt.Errorf("Not found")
	}
	index, err := l.next.Find(val)
	if err != nil {
		return -1, err
	} else {
		return 1 + index, nil
	}
}